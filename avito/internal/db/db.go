package db

import (
	"context"
	"database/sql"
	"time"
)

type dbHandler struct {
	cur *sql.DB
}

var DB = newDBConn()

func (db dbHandler) Select() {

}

func (db dbHandler) ExecSegment(users []int, segmentName string, operation bool) error {
	ctx := context.Background()
	tx, err := db.cur.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	switch operation {
	case true:
		_, err = tx.ExecContext(ctx, `INSERT INTO "userSegments"("userID", "segmentName") VALUES ($1, $2)`, users[0], segmentName)
	default:
		_, err = tx.ExecContext(ctx, `DELETE FROM "userSegments" WHERE "segmentName" = $1`, segmentName)
	}
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, id := range users {
		_, err = tx.ExecContext(ctx, `INSERT INTO "segmentHistory" ("userID", "segmentName", operation, "dateReq") VALUES ($1, $2, $3, $4)`,
			id, segmentName, operation, time.Now().Format(time.DateOnly))
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (db dbHandler) InsergUserSegments(userID int, segmentName string) error {
	_, err := db.cur.Exec(`INSERT INTO "userSegments"("userID", "segmentName") VALUES ($1, $2)`, userID, segmentName)
	if err != nil {
		return err
	}
	return nil
}
