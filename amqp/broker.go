package mq

import (
	"encoding/json"
	"errors"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	EXCHANGE_NAME = "my-exchange"
	EXCHANGE_TYPE = "topic"
)

var (
	// ErrNoActiveConn No Active Connection errors
	ErrNoActiveConn = errors.New("no Active Connection")
)

type (
	Broker interface {
		Publish() error
	}

	AMQPConfig struct{}

	AMQPBroker struct {
		sync.Mutex
		conn    *amqp.Connection
		amqpURI string
		errors  chan *amqp.Error
		channel *amqp.Channel
	}
)

func NewAMQPBroker(uri string) *AMQPBroker {
	amqpBroker := &AMQPBroker{}
	amqpBroker.amqpURI = uri

	return amqpBroker
}

func (a *AMQPBroker) Setup() error {

	conn, err := amqp.Dial(a.amqpURI)
	if err != nil {
		return err
	}

	a.Lock()
	a.conn = conn
	a.Unlock()

	errors := make(chan *amqp.Error)
	a.conn.NotifyClose(errors)
	a.errors = errors

	a.channel, _ = a.conn.Channel()

	a.declareExchange()

	return nil
}

func (a *AMQPBroker) declareExchange() {

	err := a.channel.ExchangeDeclare(
		EXCHANGE_NAME, // name
		EXCHANGE_TYPE, // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // noWait
		nil,           // arguments
	)

	if err != nil {
		log.Warn(err)
	}
}

func (a *AMQPBroker) Close() error {

	if a.conn == nil {
		return ErrNoActiveConn
	}

	return a.conn.Close()
}

func (a *AMQPBroker) Publish(payload interface{}, eventName string) error {

	body, _ := json.Marshal(payload)

	publishing := amqp.Publishing{
		Headers:      amqp.Table{},
		Body:         body,
		ContentType:  "application/json",
		DeliveryMode: amqp.Transient,
	}

	if a.conn == nil {
		return ErrNoActiveConn
	}

	channel := a.channel
	if channel == nil {
		return ErrNoActiveConn
	}

	err := channel.Publish(EXCHANGE_NAME, eventName, false, false, publishing)

	logLevel := log.InfoLevel

	if err != nil {
		logLevel = log.ErrorLevel
	}
	log.WithFields(log.Fields{
		"eventName": eventName,
	}).Info(logLevel)

	return err
}

func (a *AMQPBroker) GetConn() *amqp.Connection {

	return a.conn
}
