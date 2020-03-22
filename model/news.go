package model

import "time"

type (
	News struct {
		ID      int       `db:"id" json:"id"`
		Author  string    `db:"author" json:"author"`
		Body    string    `db:"body" json:"body"`
		Created time.Time `db:"created" json:"created"`
	}
)
