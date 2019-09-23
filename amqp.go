package easy_amqp

import "github.com/streadway/amqp"

var (
	DefaultQueueDurable       = true
	DefaultQueueAutoDelete    = false
	DefaultQueueNoWait        = false
	DefaultQueueExclusive     = false
	DefaultExchangeDurable    = true
	DefaultExchangeAutoDelete = false
	DefaultExchangeNoWait     = false
	DefaultExchangeInternal   = false
	DefaultRouteKey           = ""
	DefaultBindNoWait         = false
)

type QueueOption func(*Queue)

type ExchangeOption func(*Exchange)

type BindOption func(*Bind)

func WithBindAllRouter() BindOption {
	return func(bind *Bind) {
		bind.RouterKey = "#"
	}
}

func WithBindRouteKey(key string) BindOption {
	return func(bind *Bind) {
		bind.RouterKey = key
	}
}

func WithBindNoWait(noWait bool) BindOption {
	return func(bind *Bind) {
		bind.NoWait = noWait
	}
}

func WithBindArgs(args amqp.Table) BindOption {
	return func(bind *Bind) {
		bind.Args = args
	}
}

func WithQueueDurable(durable bool) QueueOption {
	return func(queue *Queue) {
		queue.Durable = durable
	}
}

func WithQueueAlias(alias string) QueueOption {
	return func(queue *Queue) {
		queue.Alias = alias
	}
}

func WithQueueAutoDelete(autoDelete bool) QueueOption {
	return func(queue *Queue) {
		queue.AutoDelete = autoDelete
	}
}

func WithQueueNoWait(noWait bool) QueueOption {
	return func(queue *Queue) {
		queue.NoWait = noWait
	}
}

func WithQueueArgs(args amqp.Table) QueueOption {
	return func(queue *Queue) {
		queue.Args = args
	}
}

func WithQueueExclusive(exclusive bool) QueueOption {
	return func(queue *Queue) {
		queue.Exclusive = exclusive
	}
}

func WithExchangeDurable(durable bool) ExchangeOption {
	return func(exchange *Exchange) {
		exchange.Durable = durable
	}
}

func WithExchangeAlias(alias string) ExchangeOption {
	return func(exchange *Exchange) {
		exchange.Alias = alias
	}
}

func WithExchangeDelete(autoDelete bool) ExchangeOption {
	return func(exchange *Exchange) {
		exchange.AutoDelete = autoDelete
	}
}

func WithExchangeNoWait(noWait bool) ExchangeOption {
	return func(exchange *Exchange) {
		exchange.NoWait = noWait
	}
}

func WithExchangeArgs(args amqp.Table) ExchangeOption {
	return func(exchange *Exchange) {
		exchange.Args = args
	}
}

func WithExchangeInternal(internal bool) ExchangeOption {
	return func(exchange *Exchange) {
		exchange.Internal = internal
	}
}

func NewQueue(name string, opts ...QueueOption) Queue {
	q := Queue{
		Name:       name,
		Alias:      name,
		Durable:    DefaultQueueDurable,
		AutoDelete: DefaultQueueAutoDelete,
		Exclusive:  DefaultQueueExclusive,
		NoWait:     DefaultQueueNoWait,
		Args:       amqp.Table{},
	}
	for _, o := range opts {
		o(&q)
	}
	return q
}

func NewExchange(name string, kind ExchangeKind, opts ...ExchangeOption) Exchange {
	e := Exchange{
		Name:       name,
		Kind:       kind,
		Alias:      name,
		Durable:    DefaultExchangeDurable,
		AutoDelete: DefaultExchangeAutoDelete,
		Internal:   DefaultExchangeInternal,
		NoWait:     DefaultExchangeNoWait,
		Args:       amqp.Table{},
	}
	for _, o := range opts {
		o(&e)
	}
	return e
}
