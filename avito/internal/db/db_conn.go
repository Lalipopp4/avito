package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// connection with PostgreSQL
func newDBConn() *dbHandler {
	conn := "user=postgres host=localhost password=postgres port=5432 dbname=avito sslmode=disable"
	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	return &dbHandler{
		db,
	}
}

func (db *dbHandler) Close() {
	db.Close()
}
