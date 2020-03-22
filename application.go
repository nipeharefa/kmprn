package main

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	mq "github.com/nipeharefa/kmprn/amqp"
	"github.com/nipeharefa/kmprn/controller"
	"github.com/nipeharefa/kmprn/repository"
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type (
	Application interface {
		StartHTTPServer()
		StartConsumer()
	}

	application struct {
		e             *echo.Echo
		db            *sqlx.DB
		broker        *mq.AMQPBroker
		elasticClient *elastic.Client
	}
)

// NewApplication :nodoc:
func NewApplication() Application {

	a := &application{}

	e := echo.New()
	e.Use(middleware.Recover())

	a.e = e

	a.connectToBroker()
	a.connectDB()
	a.connectToES()

	return a
}

func (a *application) connectToES() {

	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false))

	if err != nil {
		log.Fatal(err)
	}

	a.elasticClient = client
}

func (a *application) connectToBroker() {
	amqpURI := viper.GetString("amqp.uri")
	broker := mq.NewAMQPBroker(amqpURI)

	err := broker.Setup()
	if err != nil {
		log.Fatal(err)
	}

	a.broker = broker

	err = broker.Setup()
	if err != nil {
		log.Warn(err)
	}
}

func (a *application) connectDB() {

	db, err := sqlx.Connect("postgres", viper.GetString("application.db.url"))
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatalln(err)
	}

	maxidle := viper.GetInt("application.db.max_idle")
	maxConn := viper.GetInt("application.db.max_conn")

	log.Info("Database connected")
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxIdleConns(maxidle)
	db.SetMaxOpenConns(maxConn)

	a.db = db

}

func (a *application) StartHTTPServer() {

	// repository
	newsRepo := repository.NewNewsRepository(a.db)

	nc := controller.NewNewsController(a.broker, a.elasticClient, newsRepo)

	a.e.GET("/news", nc.GetNews)
	a.e.POST("/news", nc.CreateNews)

	a.e.HideBanner = true
	log.Fatal(a.e.Start(":8000"))
}
