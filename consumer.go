package main

import (
	"github.com/nipeharefa/kmprn/consumer"
	"github.com/nipeharefa/kmprn/repository"
	log "github.com/sirupsen/logrus"
)

func (a *application) StartConsumer() {

	// repository
	newsRepo := repository.NewNewsRepository(a.db)

	conn := a.broker.GetConn()
	channel, err := conn.Channel()
	if err != nil {
		log.Error(err)
	}
	newsQ := consumer.NewCreateNewsQ(channel)
	if err := newsQ.Setup(); err != nil {
		log.Error(err)
	}

	newsQHandler := consumer.NewCreateNewsHandler(newsRepo, a.elasticClient)

	newsQ.ConsumerFunc(newsQHandler.Execute)

	select {}
}
