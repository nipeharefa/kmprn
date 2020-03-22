package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nipeharefa/kmprn/model"
)

type (
	NewsRepository interface {
		Create(*model.News) error
		FindById(ID int) (model.News, error)
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

func (nr *newsRepository) FindById(ID int) (news model.News, err error) {

	query := "SELECT * FROM news where id=$1 limit 1"
	err = nr.db.Get(&news, query, ID)
	return
}
