package postgres

import (
	"database/sql"
	"fmt"
	"streamerEventViewer/cmd/config"
)

func New(config config.DB) (*sql.DB, error) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		config.User, config.Password, config.DBName)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, err
	}

	return db, nil
}
