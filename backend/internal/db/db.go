package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	database, err := sql.Open("sqlite", fmt.Sprintf("%s?_pragma=foreign_keys(ON)", path))
	if err != nil {
		return nil, err
	}

	database.SetMaxOpenConns(1)
	database.SetMaxIdleConns(1)
	database.SetConnMaxLifetime(10 * time.Minute)

	if err := database.Ping(); err != nil {
		return nil, err
	}

	return database, nil
}
