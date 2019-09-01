package easy_amqp

import (
	"github.com/streadway/amqp"
)

var (
	QueueDeclared    = make(map[string]*Queue)
	ExchangeDeclared = make(map[string]*Exchange)
)

// DeclareQueue is declare queues and bind exchange
func DeclareQueue(channel *amqp.Channel, queues ...*Queue) error {
	var exchanges []*Exchange
	for _, queue := range queues {
		q, err := channel.QueueDeclare(queue.Name, queue.Durable, queue.AutoDelete, queue.Exclusive, queue.NoWait, queue.Args)
		if err != nil {
			return err
		}
		if queue.Alias == "" {
			queue.Alias = queue.Name
		}
		queue.Queue = &q
		QueueDeclared[queue.Alias] = queue
		for _, bd := range queue.Binds {
			exchanges = append(exchanges, bd.Exg)
		}
	}
	if len(exchanges) > 0 {
		if err := DeclareExchange(channel, exchanges...); err != nil {
			return err
		}
		for _, queue := range queues {
			for _, bd := range queue.Binds {
				err := channel.QueueBind(queue.Name, bd.RouterKey, bd.Exg.Name, bd.NoWait, bd.Args)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// DeclareExchange is recursive declare exchanges and bind them
func DeclareExchange(channel *amqp.Channel, exchanges ...*Exchange) error {
	var exgs []*Exchange
	for _, exchange := range exchanges {
		err := channel.ExchangeDeclare(exchange.Name, exchange.Kind.String(), exchange.Durable, exchange.AutoDelete, exchange.Internal, exchange.NoWait, exchange.Args)
		if err != nil {
			return err
		}
		if exchange.Alias == "" {
			exchange.Alias = exchange.Name
		}
		ExchangeDeclared[exchange.Alias] = exchange
		for _, bd := range exchange.Binds {
			exgs = append(exgs, bd.Exg)
		}
	}
	if len(exgs) > 0 {
		if err := DeclareExchange(channel, exgs...); err != nil {
			return err
		}
		for _, exchange := range exchanges {
			for _, bd := range exchange.Binds {
				err := channel.ExchangeBind(exchange.Name, bd.RouterKey, bd.Exg.Name, bd.NoWait, bd.Args)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
