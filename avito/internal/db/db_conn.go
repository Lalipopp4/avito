package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewDBConn() *sql.DB {
	conn := "user=postgres host=localhost password=postgres port=5432 dbname=avito sslmode=disable"
	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	return db
}

func Close(db *sql.DB) {
	db.Close()
}
