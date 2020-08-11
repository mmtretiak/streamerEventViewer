package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"streamerEventViewer/cmd/config"
)

func New(config config.DB) (*sql.DB, error) {
	dbURL := getDatabaseURL(config)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getDatabaseURL(config config.DB) string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}

	dbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	return dbInfo
}
