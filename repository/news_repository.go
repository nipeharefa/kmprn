package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nipeharefa/kmprn/model"
)

type (
	NewsRepository interface {
		Create(*model.News) error
	}
	newsRepository struct {
		db *sqlx.DB
	}
)

func NewNewsRepository(db *sqlx.DB) NewsRepository {

	nr := &newsRepository{db}

	return nr
}

func (nr *newsRepository) Create(news *model.News) error {

	news.Created = time.Now()

	query := `INSERT INTO news(author, body, created) values($1, $2, $3) returning id`

	err := nr.db.QueryRow(query, news.Author, news.Body, news.Created).Scan(&news.ID)
	if err != nil {
		return err
	}

	return nil
}
