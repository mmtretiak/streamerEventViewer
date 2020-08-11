package postgres

import (
	"database/sql"
	"os"
	"streamerEventViewer/cmd/config"
)

func New(config config.DB) (*sql.DB, error) {
	//dbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	//	"password=%s dbname=%s sslmode=disable",
	//	config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	return db, nil
}
