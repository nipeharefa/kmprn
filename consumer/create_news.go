package consumer

import (
	"fmt"

	"github.com/streadway/amqp"
)

type (
	createNewsQ struct {
		channel *amqp.Channel
	}
)

func NewCreateNewsQ(c *amqp.Channel) Queue {

	cq := createNewsQ{}
	cq.channel = c

	return cq
}

func (s createNewsQ) GetBindKey() string {
	return "news.created"
}

func (s createNewsQ) GetQueueName() string {
	return "news_created"
}

func (s createNewsQ) Consume() (<-chan amqp.Delivery, error) {
	deliveries, err := s.channel.Consume(
		s.GetQueueName(), // name
		"",               // consumerTag,
		false,            // noAck
		false,            // exclusive
		false,            // noLocal
		false,            // noWait
		nil,              // arguments
	)

	return deliveries, err
}

func (s createNewsQ) Setup() error {

	// Declare
	_, err := s.channel.QueueDeclare(
		s.GetQueueName(), // name of the queue
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // noWait
		nil,              // arguments
	)

	if err != nil {
		return err
	}

	err = s.channel.QueueBind(
		s.GetQueueName(), // name of the queue
		"news.created",   // bindingKey
		"my-exchange",    // sourceExchange
		false,            // noWait
		nil,              // arguments
	)

	if err != nil {
		return err
	}

	return nil
}

func (s createNewsQ) ConsumerFunc(fn ConsumerFunc) error {
	deliveries, err := s.Consume()
	if err != nil {
		return err
	}

	var worker = func(id int, jobs <-chan amqp.Delivery) {
		for j := range jobs {
			fn(j)
			fmt.Println("worker", id)
		}
	}

	for w := 1; w <= 5; w++ {
		go worker(w, deliveries)
	}

	return nil
}
