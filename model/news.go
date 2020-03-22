package model

import "time"

type (
	News struct {
		ID      int       `db:"id"`
		Author  string    `db:"author"`
		Body    string    `db:"body"`
		Created time.Time `db:"created"`
	}
)
