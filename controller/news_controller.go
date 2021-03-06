package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	mq "github.com/nipeharefa/kmprn/amqp"
	"github.com/nipeharefa/kmprn/model"
	"github.com/nipeharefa/kmprn/repository"
	"github.com/olivere/elastic/v7"
	"github.com/patrickmn/go-cache"
)

type (
	// NewsController interface :nodoc:
	NewsController interface {
		CreateNews(c echo.Context) error
		GetNews(c echo.Context) error
	}

	newsController struct {
		broker      *mq.AMQPBroker
		esClient    *elastic.Client
		newsRepo    repository.NewsRepository
		cacheSystem *cache.Cache
	}

	createNewsRequest struct {
		Author string `json:"author"`
		Body   string `json:"body"`
	}

	esSearcHResult struct {
		ID      int       `json:"id"`
		Created time.Time `json:"created"`
	}
)

// NewNewsController :nodoc:
func NewNewsController(broker *mq.AMQPBroker, esClient *elastic.Client, newsRepo repository.NewsRepository) NewsController {

	nc := &newsController{}
	nc.broker = broker
	nc.esClient = esClient
	nc.newsRepo = newsRepo

	nc.cacheSystem = cache.New(5*time.Minute, 10*time.Minute)

	return nc
}

func (nc *newsController) GetNews(c echo.Context) error {

	page := 0

	page, _ = strconv.Atoi(c.QueryParam("page"))

	n := nc.search(nc.esClient, page)
	return c.JSON(http.StatusOK, n)
}

func (nc *newsController) CreateNews(c echo.Context) error {

	req := createNewsRequest{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	_ = nc.broker.Publish(req, "news.created")

	return c.JSON(http.StatusCreated, nil)
}

func (nc *newsController) search(esClient *elastic.Client, page int) []model.News {

	var wg sync.WaitGroup

	newses := make([]model.News, 0)
	perPage := 10

	ctx := context.Background()

	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchAllQuery())
	searchSource.Size(perPage)
	searchSource.From(perPage * page)
	searchSource.Sort("created", false)

	searchService := esClient.Search().
		Index("news").
		SearchSource(searchSource)

	searchResult, err := searchService.Do(ctx)
	if err != nil {
		fmt.Println(err)
		return newses
	}

	totalHits := len(searchResult.Hits.Hits)

	if totalHits == 0 {
		return newses
	}

	var firstRow esSearcHResult

	firstRecord := searchResult.Hits.Hits[0]
	_ = json.Unmarshal(firstRecord.Source, &firstRow)

	newsInCache, found := nc.cacheSystem.Get(string(firstRow.ID))
	if found {
		return newsInCache.([]model.News)
	}

	newsChan := make(chan model.News, totalHits)

	wg.Add(totalHits)
	for _, hit := range searchResult.Hits.Hits {

		var news esSearcHResult
		_ = json.Unmarshal(hit.Source, &news)

		go nc.findNews(newsChan, &wg, news.ID)
	}

	wg.Wait()
	close(newsChan)

	for n := range newsChan {
		newses = append(newses, n)
	}

	// Sort News By ID
	sort.Slice(newses, func(i, j int) bool {
		return newses[i].ID > newses[j].ID
	})

	nc.cacheSystem.Set(string(newses[0].ID), newses, cache.DefaultExpiration)

	return newses
}

func (nc *newsController) findNews(ch chan model.News, wg *sync.WaitGroup, id int) {

	defer wg.Done()

	a, err := nc.newsRepo.FindByID(id)
	if err != nil {
		return
	}

	ch <- a
}
