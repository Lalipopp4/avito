package service

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Lalipopp4/avito/internal/logger"
)

type point struct {
	userID      int
	segmentName string
}

var (
	// map of added segments
	segments = make(map[string]struct{})
	// map of temporary segments for user
	ttlUSers = make(map[time.Time][]point)
)

// function to decode data in request
func decode(d interface{}, w http.ResponseWriter, r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		logger.Logger.Log(err)
		w.Write([]byte(error.Error(err)))
		return err
	}
	return nil
}

// adding segments
func addSegment(w http.ResponseWriter, r *http.Request) {
	var segment segment
	err := decode(&segment, w, r)
	if err != nil {
		logger.Logger.Log(err)
		return
	}

	//checking if segment already exists
	res, err := DB.Exec("SELECT * FROM segments WHERE name = $1", segment.Name)

	if err != nil {
		logger.Logger.Log(err)
		w.Write([]byte("Error: error in SQL request.\n"))
		return
	}
	if n, err := res.RowsAffected(); n > 0 || string(segment.Name) == "" || err != nil {
		logger.Logger.Log("Error: this segment is already exists or empty request.\n")
		w.Write([]byte("Error: this segment is already exists or empty request.\n"))
		return
	}

	ctx := context.Background()
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		logger.Logger.Log(err)
		w.Write([]byte("Error: error in SQL request.\n"))
		return
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO segments VALUES ($1)", segment.Name)
	if err != nil {
		logger.Logger.Log(err)
		tx.Rollback()
		w.Write([]byte("Error: error in SQL request.\n"))
		return
	}

	res, err = DB.Exec("SELECT * FROM users")
	n, err := res.RowsAffected()
	for i := 0; i < int(n); i += int(n / int64(segment.Perc)) {
		_, err = tx.ExecContext(ctx, `INSERT INTO "segmentHistory" ("userID", "segmentName", operation, "dateReq") VALUES ($1, $2, $3, $4)`,
			i, segment.Name, 1, time.Now().Format(time.DateOnly))
		if err != nil {
			logger.Logger.Log(err)
			tx.Rollback()
			w.Write([]byte("Error: error in SQL request.\n"))
			return
		}

		_, err = tx.ExecContext(ctx, `INSERT INTO "userSegments"("userID", "segmentName") VALUES ($1, $2)`, i, segment.Name)
		if err != nil {
			logger.Logger.Log(err)
			tx.Rollback()
			w.Write([]byte("Error: error in SQL request.\n"))
			return
		}

		err = tx.Commit()
		if err != nil {
			logger.Logger.Log(err)
			w.Write([]byte("Error: error in SQL request.\n"))
			return
		}
	}

	logger.Logger.Log(segment.Name + " added.\n")
	w.Write([]byte(segment.Name + " added.\n"))
}

// deleting segment
func deleteSegment(w http.ResponseWriter, r *http.Request) {
	var segment segment
	err := decode(&segment, w, r)
	if err != nil {
		return
	}

	//checking if segment doesn't exist
	fmt.Println(segments)
	if _, ok := segments[string(segment.Name)]; !ok || string(segment.Name) == "" {
		logger.Logger.Log("Error: this segment doesn't exist or empty request.\n")
		w.Write([]byte("Error: this segment doesn't exist or empty request.\n"))
		return
	}

	// SQL transaction to find users in this segment,
	// delete segment and insert info about deleting users from segment
	// ------
	ctx := context.Background()
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		logger.Logger.Log(err)
		return
	}
	rows, err := tx.QueryContext(ctx, `SELECT "userID" FROM "userSegments" WHERE "segmentName" = $1`, segment.Name)
	if err != nil {
		logger.Logger.Log(err)
		tx.Rollback()
		w.Write([]byte("Error: error in SQL request.\n"))
		return
	}
	var (
		id    int
		users []int
	)
	for rows.Next() {
		rows.Scan(&id)
		users = append(users, id)
	}
	_, err = tx.ExecContext(ctx, `DELETE FROM "userSegments" WHERE "segmentName" = $1`, segment.Name)
	if err != nil {
		logger.Logger.Log(err)
		tx.Rollback()
		w.Write([]byte("Error: error in SQL request.\n"))
		return
	}
	for _, val := range users {
		_, err = tx.ExecContext(ctx, `INSERT INTO "segmentHistory" ("userID", "segmentName", operation, "dateReq") VALUES ($1, $2, $3, $4)`,
			val, segment.Name, 0, time.Now().Format(time.DateOnly))
		if err != nil {
			logger.Logger.Log(err)
			tx.Rollback()
			w.Write([]byte("Error: error in SQL request.\n"))
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		logger.Logger.Log(err)
		w.Write([]byte("Error: error in SQL request.\n"))
		return
	}
	// ------

	delete(segments, segment.Name)
	logger.Logger.Log(segment.Name + " deleted.\n")
	w.Write([]byte(segment.Name + " deleted.\n"))
}

// adding user, adding and deleting segments for him
func addUser(w http.ResponseWriter, r *http.Request) {
	var request addUserRequest
	err := decode(&request, w, r)
	if err != nil {
		return
	}

	// checking if segments from add list doesn't exist
	for _, segment := range append(request.SegmentsToAdd, request.SegmetsToDelete...) {
		if _, ok := segments[segment]; !ok {
			logger.Logger.Log("Error: no needed segment.")
			w.Write([]byte("Error: no such segments"))
			return
		}
	}

	// checking if there is intersection of adding and deleting segments
	intersection := make(map[string]struct{})
	for _, segment := range request.SegmentsToAdd {
		intersection[segment] = struct{}{}
	}
	for _, val := range request.SegmetsToDelete {
		if _, ok := intersection[val]; ok {
			logger.Logger.Log(err)
			w.Write([]byte("Error: common segments in add ad delete lists.\n"))
			return
		}
	}

	// adding user in segments
	for _, val := range request.SegmentsToAdd {
		var n int

		// checking if user is already in segment
		rows, err := DB.Query(`SELECT COUNT(*) FROM "userSegments" WHERE "segmentName" = $1`, val)
		rows.Next()
		if rows.Scan(&n); n > 0 || err != nil {

			continue
		}

		// SQL transaction to insert data in segmentHistory and userSegment
		// ------
		ctx := context.Background()
		tx, err := DB.BeginTx(ctx, nil)
		if err != nil {
			logger.Logger.Log(err)
			w.Write([]byte("Error: error in SQL request.\n"))
			return
		}

		_, err = tx.ExecContext(ctx, `INSERT INTO "segmentHistory" ("userID", "segmentName", operation, "dateReq") VALUES ($1, $2, $3, $4)`,
			request.UserID, val, 1, time.Now().Format(time.DateOnly))
		if err != nil {
			logger.Logger.Log(err)
			tx.Rollback()
			w.Write([]byte("Error: error in SQL request.\n"))
			return
		}

		_, err = tx.ExecContext(ctx, `INSERT INTO "userSegments"("userID", "segmentName") VALUES ($1, $2)`, request.UserID, val)
		if err != nil {
			logger.Logger.Log(err)
			tx.Rollback()
			w.Write([]byte("Error: error in SQL request.\n"))
			return
		}

		err = tx.Commit()
		if err != nil {
			logger.Logger.Log(err)
			w.Write([]byte("Error: error in SQL request.\n"))
			return
		}

		//ttlUSers[time.Now()+time.Day()]
		// ------
	}
	logger.Logger.Log("Segments added.\n")
	w.Write([]byte("Segments added.\n"))

	// deleting user from segments
	for _, val := range request.SegmetsToDelete {
		var n int

		// checking if user is already in segment
		rows, err := DB.Query(`SELECT COUNT(*) FROM "userSegments" WHERE "segmentName" = $1`, val)
		if rows.Scan(&n); n == 0 || err != nil {
			continue
		}

		// SQL transaction to insert data in segmentHistory and userSegment
		// ------
		ctx := context.Background()
		tx, err := DB.BeginTx(ctx, nil)
		if err != nil {
			logger.Logger.Log(err)
			return
		}

		_, err = tx.ExecContext(ctx, `INSERT INTO "segmentHistory" ("userID", "segmentName", operation, "dateReq") VALUES ($1, $2, $3, $4)`,
			request.UserID, val, 0, time.Now().Format(time.DateOnly))
		if err != nil {
			logger.Logger.Log(err)
			tx.Rollback()
			w.Write([]byte("Error: error in SQL request.\n"))
			return
		}

		_, err = tx.ExecContext(ctx, `DELETE FROM "userSegments" WHERE "segmentName" = $1`, val)
		if err != nil {
			logger.Logger.Log(err)
			tx.Rollback()
			w.Write([]byte("Error: error in SQL request.\n"))
			return
		}

		err = tx.Commit()
		if err != nil {
			logger.Logger.Log(err)
			w.Write([]byte("Error: error in SQL request.\n"))
			return
		}

		// ------
	}

	logger.Logger.Log("Segments deleted.\n")
	w.Write([]byte("Segments deleted.\n"))
}

// searching for active segments of user
func activeUserSegments(w http.ResponseWriter, r *http.Request) {
	var userID userRequest
	err := decode(&userID, w, r)
	if err != nil {
		return
	}

	rows, err := DB.Query(`SELECT "segmentName" FROM "userSegments" WHERE "userID" = $1`, userID.UserID)
	if err != nil {
		w.Write([]byte("Error: error in SQL request.\n"))
		logger.Logger.Log(err)
		return
	}
	var (
		userSegments []string
		segment      string
	)
	for rows.Next() {
		rows.Scan(&segment)
		userSegments = append(userSegments, segment)
	}
	response, err := json.Marshal(userSegments)
	if err != nil {
		w.Write([]byte("Error: error in json encoding.\n"))
		logger.Logger.Log(err)
		return
	}

	logger.Logger.Log(response)
	w.Write(response)
}

func segmentHistory(w http.ResponseWriter, r *http.Request) {
	var historyReq historyRequest
	err := decode(&historyReq, w, r)
	if err != nil {
		return
	}
	rows, err := DB.Query(`SELECT "userID", "SegmentName", operation, "dateReq" FROM "segmentHistory"
							WHERE "dateReq" = $1`, historyReq.Date)
	if err != nil {
		w.Write([]byte("Error: error in SQL request.\n"))
		logger.Logger.Log(err)
		return
	}
	f, err := os.Create("github.com/Lalipopp4/avito/pkg/files/history.csv")
	defer f.Close()

	if err != nil {
		w.Write([]byte("Error: error in openning file.\n"))
		logger.Logger.Log(err)
		return
	}
	var (
		id              int
		name, date, opS string
		op              bool
	)
	csvW := csv.NewWriter(f)
	defer csvW.Flush()
	for rows.Next() {
		rows.Scan(&id, &name, &op, &date)
		if op {
			opS = "Adding"
		} else {
			opS = "Deleting"
		}
		if err := csvW.Write([]string{strconv.Itoa(id), name, opS, date}); err != nil {
			w.Write([]byte("Error: error with file.\n"))
			logger.Logger.Log(err)
			return
		}
	}

}
