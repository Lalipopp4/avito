package db

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dbHandler struct {
	cur *gorm.DB
}

var DB = newDBConn()

// conn function for DB
func newDBConn() *dbHandler {
	conn := "user=postgres host=localhost password=postgres port=5432 dbname=avito sslmode=disable"
	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &dbHandler{
		db,
	}
}

type historyReq struct {
	UserID        int
	SegmentName   string
	Operation     bool
	Timestamp_req time.Time
}

type segReq struct {
	UserID      int
	SegmentName string
	TTL         string
}
