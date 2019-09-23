package easy_amqp_test

import (
	"context"
	easy_amqp "github.com/speeder-allen/easy-amqp"
	"github.com/streadway/amqp"
	"gotest.tools/assert"
	"log"
	"testing"
	"time"
)

func TestNewSubscriber(t *testing.T) {
	sub := easy_amqp.NewSubscriber(nil, context.Background())
	sub.Subscribe("test", easy_amqp.NewSUbscribeDefaultOption(func(delivery *amqp.Delivery) {
		log.Println(delivery)
	}))
	err := sub.Comsume(2)
	assert.Equal(t, err, easy_amqp.ErrorConnection)
}

func TestSubscriber_Comsume(t *testing.T) {
	conn, err := amqp.Dial("amqp://root:123456@localhost:5672")
	if err != nil {
		t.Fatalf("dail rabbitmq faild %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	sub := easy_amqp.NewSubscriber(conn, ctx)
	sub.Subscribe("test1", easy_amqp.NewSUbscribeDefaultOption(func(delivery *amqp.Delivery) {
		log.Println(delivery)
	}))
	sub.Subscribe("test2", easy_amqp.NewSUbscribeDefaultOption(func(delivery *amqp.Delivery) {
		log.Println(delivery)
	}))
	sub.Subscribe("test3", easy_amqp.NewSUbscribeDefaultOption(func(delivery *amqp.Delivery) {
		log.Println(delivery)
	}))
	sub.Comsume(3)
}
