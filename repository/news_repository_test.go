package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/nipeharefa/kmprn/model"
	"github.com/stretchr/testify/assert"
)

func TestNewsRepository(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer mockDB.Close()

	t.Run("TestFindByID", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "author", "body", "created"}).
			AddRow(1, "kumparan", "lorem lorem", time.Now())

		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT * FROM news where id=$1 limit 1")).
			WillReturnRows(rows)

		db := sqlx.NewDb(mockDB, "sqlmock")
		newsRepo := NewNewsRepository(db)

		news, err := newsRepo.FindByID(1)
		assert.NoError(t, err)
		assert.Equal(t, news.Author, "kumparan")
	})

	t.Run("TestErrFindByID", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "author", "body", "created"})

		mock.ExpectQuery(
			regexp.QuoteMeta("SELECT * FROM news where id=$1 limit 1")).
			WillReturnRows(rows)

		db := sqlx.NewDb(mockDB, "sqlmock")
		newsRepo := NewNewsRepository(db)

		_, err := newsRepo.FindByID(2)
		assert.Error(t, err, sql.ErrNoRows)
	})

	t.Run("TestCreate", func(t *testing.T) {

		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(1)

		query := "INSERT INTO news(author, body, created) values($1, $2, $3) returning id"

		mock.ExpectQuery(
			regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		news := model.News{}
		news.Author = "kumparan"
		news.Body = "kumparan content"

		db := sqlx.NewDb(mockDB, "sqlmock")
		newsRepo := NewNewsRepository(db)

		err := newsRepo.Create(&news)
		assert.NoError(t, err)

	})

	t.Run("TestErrCreate", func(t *testing.T) {

		sampleError := errors.New("something wrong")
		query := "INSERT INTO news(author, body, created) values($1, $2, $3) returning id"

		mock.ExpectQuery(
			regexp.QuoteMeta(query)).
			WillReturnError(sampleError)

		news := model.News{}
		news.Author = "kumparan"
		news.Body = "kumparan content"

		db := sqlx.NewDb(mockDB, "sqlmock")
		newsRepo := NewNewsRepository(db)

		err := newsRepo.Create(&news)
		assert.Error(t, err, sampleError)

	})
}
