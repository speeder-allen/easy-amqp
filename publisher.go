package easy_amqp

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrorClosed         = errors.New("publisher is closed")
	ErrorAmqpConnection = errors.New("error amqp connection")
	ErrorNoChannel      = errors.New("not channel gets")
	ErrorBuildChannel   = errors.New("can't build channel")
)

type publisher struct {
	conn       *amqp.Connection
	ctx        context.Context
	uuid       string
	cancelFunc context.CancelFunc
	sync.Once
	closed           uint32
	defaultExchange  string
	defaultRouterKey string
	pool             *channelPool
}

type channelFactory func() *amqp.Channel

type channelPool struct {
	capacity uint
	channels chan *amqp.Channel
	factory  channelFactory
	rw       sync.RWMutex
}

func (c *channelPool) Len() int {
	return int(c.capacity)
}

func (c *channelPool) Take() (*amqp.Channel, error) {
	c.rw.RLock()
	cns := c.channels
	c.rw.RUnlock()
	select {
	case cnl := <-cns:
		if cnl == nil {
			return nil, ErrorNoChannel
		}
		return cnl, nil
	default:
		cnl := c.factory()
		if cnl == nil {
			return nil, ErrorBuildChannel
		}
		return cnl, nil
	}
}

func (c *channelPool) Put(channel *amqp.Channel) error {
	c.rw.RLock()
	defer c.rw.RUnlock()
	select {
	case c.channels <- channel:
		return nil
	default:
		return channel.Close()
	}
}

func NewChannelPool(capacity uint, factory channelFactory) *channelPool {
	p := channelPool{
		capacity: capacity,
		channels: make(chan *amqp.Channel, capacity),
		factory:  factory,
	}
	return &p
}

func (p *publisher) SetDefaultExchange(exchange string) {
	p.defaultExchange = exchange
}

func (p *publisher) SetDefaultRouterKey(routerKey string) {
	p.defaultRouterKey = routerKey
}

func (p *publisher) Publishe(msg *message) error {
	if atomic.LoadUint32(&p.closed) == 1 {
		return ErrorClosed
	}
	if p.conn == nil {
		return ErrorAmqpConnection
	}
	channel, err := p.pool.Take()
	if err != nil {
		return err
	}
	if msg.Exchange == "" {
		msg.Exchange = p.defaultExchange
	}
	if msg.RouterKey == "" {
		msg.RouterKey = p.defaultRouterKey
	}
	bytes, err := msg.Body()
	if err != nil {
		return err
	}
	return channel.Publish(msg.Exchange, msg.RouterKey, true, false, amqp.Publishing{
		ContentType:  msg.ContentType.String(),
		Body:         bytes,
		AppId:        msg.AppId,
		MessageId:    msg.UUID,
		Type:         msg.Type,
		DeliveryMode: msg.Mode,
		Priority:     msg.Priority,
		Timestamp:    time.Now(),
	})
}

func (p *publisher) Close() error {
	if atomic.CompareAndSwapUint32(&p.closed, 0, 1) {
		p.cancelFunc()
	}
	return nil
}

func NewPublisher(conn *amqp.Connection, ctx context.Context) *publisher {
	ctx, cancel := context.WithCancel(ctx)
	return &publisher{
		conn:       conn,
		ctx:        ctx,
		uuid:       uuid.New().String(),
		cancelFunc: cancel,
		pool: NewChannelPool(5, func() *amqp.Channel {
			if channel, err := conn.Channel(); err == nil {
				return channel
			}
			return nil
		}),
	}
}
