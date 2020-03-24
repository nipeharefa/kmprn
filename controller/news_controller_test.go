package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/example/federation/accounts/graph/model"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	mq "github.com/nipeharefa/kmprn/amqp"
	"github.com/nipeharefa/kmprn/repository"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
)

func TestNewsController(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repository.NewMockNewsRepository(ctrl)
	handler := NewNewsController(&mq.AMQPBroker{}, &elastic.Client{}, repo)

	t.Run("TestBadRequestCreateNews", func(t *testing.T) {

		userJSON := `{`

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, handler.CreateNews(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("TestcreateNews", func(t *testing.T) {

		newsJSON := `{"author": "Kumparan","body": "Hands On Oppo Reno 3 yang Rilis di Indonesia, Jagokan 4 Kamera Belakang"}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(newsJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		sampleNews := model.User{}

		repo.EXPECT().Create(&sampleNews).Return(nil).AnyTimes()

		if assert.NoError(t, handler.CreateNews(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	})

}
