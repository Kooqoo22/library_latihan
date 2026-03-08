package config

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectDB() (*sqlx.DB, error) {

	dsn := "postgres://faizrabbani:password@localhost:5432/library_practice_db?sslmode=disable"

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	log.Println("database connected")

	return db, nil
}