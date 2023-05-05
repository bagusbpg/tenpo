package sql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Config struct {
	Username           string
	Password           string
	Host               string
	DBName             string
	SSL                string
	MaxOpenConnections int
	MaxIdleConnections int
}

func NewClient(config Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", config.Username, config.Password, config.Host, config.DBName, config.SSL))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %s", err.Error())
	}

	db.SetMaxOpenConns(config.MaxOpenConnections)
	db.SetMaxIdleConns(config.MaxIdleConnections)

	return db, nil
}
