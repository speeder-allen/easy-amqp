package easy_amqp

import (
	"github.com/streadway/amqp"
	"strconv"
	"sync"
)

// Queue is define amqp queue
type Queue struct {
	Name       string
	Alias      string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
	Queue      *amqp.Queue
	sync.Mutex
	Binds []Bind
}

func (q *Queue) Bind(ex Exchange, opts ...BindOption) {
	bind := Bind{Exg: ex, RouterKey: DefaultRouteKey, NoWait: DefaultBindNoWait, Args: amqp.Table{}}
	for _, o := range opts {
		o(&bind)
	}
	q.Lock()
	defer q.Unlock()
	q.Binds = append(q.Binds, bind)
}

// Exchange is define amqp exchange
type Exchange struct {
	Name       string
	Alias      string
	Kind       ExchangeKind
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
	sync.Mutex
	Binds []*Bind
}

func (e *Exchange) Bind(ex Exchange, opts ...BindOption) {
	bind := Bind{Exg: ex, RouterKey: DefaultRouteKey, NoWait: DefaultBindNoWait, Args: amqp.Table{}}
	for _, o := range opts {
		o(&bind)
	}
	e.Lock()
	defer e.Unlock()
	e.Binds = append(e.Binds, &bind)
}

// Bind is define exchange relation to queue or exchange relation to exchange
type Bind struct {
	Exg       Exchange
	RouterKey string
	NoWait    bool
	Args      amqp.Table
}

// ExchangeKind
type ExchangeKind uint8

const (
	Direct ExchangeKind = iota
	Fanout
	Topic
	Headers
)

var kinds = map[ExchangeKind]string{0: amqp.ExchangeDirect, 1: amqp.ExchangeFanout, 2: amqp.ExchangeTopic, 3: amqp.ExchangeHeaders}

func (e ExchangeKind) String() string {
	if s, ok := kinds[e]; ok {
		return s
	}
	return strconv.FormatUint(uint64(e), 10)
}
