package repository

import (
	"database/sql"
	"regexp"
)

type repository struct {
	db *sql.DB
}

var whitespaceNormalizer = regexp.MustCompile(`\s+`)

func New(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}
