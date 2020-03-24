package consumer

import (
	"github.com/streadway/amqp"
)

type (
	// Queue interface
	Queue interface {
		GetQueueName() string
		Consume() (<-chan amqp.Delivery, error)
		ConsumerFunc(fn Handler) error
		Setup() error
	}

	// QueueHandler interfce
	QueueHandler interface {
		Execute(amqp.Delivery)
	}

	// Handler :nodoc:
	Handler func(amqp.Delivery)
)
