package easy_amqp

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/speeder-allen/easy-amqp/pool"
	"github.com/streadway/amqp"
	"reflect"
	"sync"
	"sync/atomic"
)

var (
	defaultConcurrencyCount  int64 = 10
	ErrorConnection                = errors.New("error connection of amqp")
	ErrorChanRecvFail              = errors.New("channel recv fail")
	ErrorEmptySubscribeQueue       = errors.New("subscribe queues empty")
)

type Handle func(delivery *amqp.Delivery)

type subscriber struct {
	conn       *amqp.Connection
	ctx        context.Context
	uuid       string
	mu         sync.Mutex
	queues     map[string]*SubscribeOption
	cancelFunc context.CancelFunc
	closed     int32
}

func (s *subscriber) Subscribe(queuename string, option *SubscribeOption) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queues[queuename] = option
}

func (s *subscriber) Comsume(concurrency int64) error {
	if s.conn == nil {
		return ErrorConnection
	}
	if len(s.queues) == 0 {
		return ErrorEmptySubscribeQueue
	}
	channel, err := s.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	if concurrency <= 0 {
		concurrency = defaultConcurrencyCount
	}
	p := pool.NewPool(concurrency)
	defer p.Close()
	var fns []Handle
	var consumes []interface{}
	consumes = append(consumes, s.ctx.Done())
	for na, opt := range s.queues {
		ch, err := channel.Consume(na, s.uuid+"-"+na, opt.AutoAck, opt.Exclusive, opt.NoLocal, opt.NoWait, opt.Args)
		if err != nil {
			return err
		}
		consumes = append(consumes, ch)
		fns = append(fns, opt.Hande)
	}
	cases := make([]reflect.SelectCase, len(consumes))
	for idx := 0; idx < len(consumes); idx++ {
		cases[idx] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(consumes[idx])}
	}
loop:
	for {
		chosen, recv, recvOk := reflect.Select(cases)
		if chosen == 0 {
			break loop
		}
		if recvOk {
			hande := fns[chosen-1]
			delivery := recv.Interface().(amqp.Delivery)
			p.Get()
			go func() {
				defer p.Put()
				hande(&delivery)
			}()
		} else {
			return ErrorChanRecvFail
		}
	}
	return nil
}

func (s *subscriber) Close() {
	if atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		s.cancelFunc()
	}
}

type SubscribeOption struct {
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
	Hande     Handle
}

func NewSUbscribeDefaultOption(handle Handle) *SubscribeOption {
	return &SubscribeOption{
		AutoAck:   false,
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
		Args:      amqp.Table{},
		Hande:     handle,
	}
}

func NewSubscriber(conn *amqp.Connection, ctx context.Context) *subscriber {
	ctx, cancel := context.WithCancel(ctx)
	return &subscriber{
		conn:       conn,
		ctx:        ctx,
		cancelFunc: cancel,
		queues:     make(map[string]*SubscribeOption),
		uuid:       uuid.New().String(),
	}
}
