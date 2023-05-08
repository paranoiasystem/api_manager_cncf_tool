package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func OpenConnection() (*sql.DB, error) {
	connStr := "user=postgres password=example dbname=postgres host=localhost port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}
