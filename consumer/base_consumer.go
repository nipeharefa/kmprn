package consumer

import (
	"github.com/streadway/amqp"
)

type (
	Queue interface {
		GetQueueName() string
		Consume() (<-chan amqp.Delivery, error)
		ConsumerFunc(fn ConsumerFunc) error
		Setup() error
	}

	QueueHandler interface {
		Execute(amqp.Delivery)
	}

	ConsumerFunc func(amqp.Delivery)
)
