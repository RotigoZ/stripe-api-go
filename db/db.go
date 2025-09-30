package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")

	db, erro := sql.Open("postgres", connStr)
	if erro != nil {
		log.Fatalf("Error oppening the connection with the database: %v", erro)
		return nil, erro
	}

	if erro := db.Ping(); erro != nil {
		log.Fatalf("Error connecting to the database: %v", erro)
		return nil, erro
	}
	return db, nil
}
