package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	mq "github.com/nipeharefa/kmprn/amqp"
)

type (
	NewsController interface {
		CreateNews(c echo.Context) error
	}

	newsController struct {
		broker *mq.AMQPBroker
	}

	createNewsRequest struct {
		Author string `json:"author"`
		Body   string `json:"body"`
	}
)

func NewNewsController(broker *mq.AMQPBroker) NewsController {

	nc := &newsController{}
	nc.broker = broker

	return nc
}

func (nc *newsController) CreateNews(c echo.Context) error {

	req := createNewsRequest{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	nc.broker.Publish(req, "news.created")

	return c.JSON(http.StatusCreated, nil)
}
