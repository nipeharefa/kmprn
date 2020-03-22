package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nipeharefa/kmprn/model"
	"github.com/nipeharefa/kmprn/repository"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type (
	createNewsHandler struct {
		newsRepo      repository.NewsRepository
		elasticClient *elastic.Client
	}

	newsData struct {
		Author string `json:"author"`
		Body   string `json:"body"`
	}

	esModel struct {
		ID      int       `json:"id"`
		Created time.Time `json:"created"`
	}
)

func NewCreateNewsHandler(newsRepo repository.NewsRepository, elasticClient *elastic.Client) QueueHandler {

	return &createNewsHandler{newsRepo, elasticClient}
}

func (h *createNewsHandler) Execute(d amqp.Delivery) {
	data := newsData{}
	err := json.Unmarshal(d.Body, &data)
	if err != nil {
		return
	}

	news := model.News{}
	news.Author = data.Author
	news.Body = data.Body

	err = h.newsRepo.Create(&news)
	if err != nil {
		fmt.Println(err)
	}

	esData := esModel{}
	esData.ID = news.ID
	esData.Created = news.Created

	esDataJSON, err := json.Marshal(esData)
	jsonString := string(esDataJSON)

	ctx := context.Background()
	_, err = h.elasticClient.Index().
		Index("news").
		BodyJson(jsonString).
		Do(ctx)

	if err != nil {
		logrus.Error(err)
	}

	_ = d.Ack(false)
}
