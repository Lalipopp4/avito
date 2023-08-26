package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Lalipopp4/avito/internal/db"
)

var (
	// sql handler
	DB = db.NewDBConn()
	// map of added segments
	segments = make(map[string]struct{})
)

type segment struct {
	Name string `json:"segmentName"`
	Perc int    `json:"percent"`
}

type addUserRequest struct {
	SegmentsToAdd   []string `json:"segmentsToAdd"`
	SegmetsToDelete []string `json:"segmentsToDelete"`
	UserID          int      `json:"userID"`
}

type userRequest struct {
	SegmentName string `json:"segmentName"`
	UserID      int    `json:"userID"`
}

// function to decode data in request
func decode(d interface{}, w http.ResponseWriter, r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return
	}

	//checking if segment already exists
	if _, ok := segments[string(segment.Name)]; ok || string(segment.Name) == "" {
		w.Write([]byte("Error: this segment is already exists or empty request.\n"))
		return
	}

	//logger.Logger.Log("POST", sg)
	segments[segment.Name] = struct{}{}
	log.Println(segment.Name + " added.")
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
		w.Write([]byte("Error: this segment doesn't exist or empty request.\n"))
		return
	}

	// SQL transaction to find users in this segment,
	// delete segment and insert info about deleting users from segment
	// ------
	ctx := context.Background()
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return
	}
	rows, err := tx.QueryContext(ctx, `SELECT "userID" FROM "userSegments" WHERE "segmentName" = $1`, segment.Name)
	if err != nil {
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
		tx.Rollback()
		w.Write([]byte("Error: error in SQL request.\n"))
		return
	}
	for _, val := range users {
		_, err = tx.ExecContext(ctx, `INSERT INTO "segmentHistory" ("userID", "segmentName", operation, "dateReq") VALUES ($1, $2, $3, $4)`,
			val, segment.Name, 0, time.Now().Format(time.DateOnly))
		if err != nil {
			tx.Rollback()
			w.Write([]byte("Error: error in SQL request.\n"))
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error: error in SQL request.\n"))
		return
	}
	// ------

	delete(segments, segment.Name)
	//logger.Logger.Log("POST", sg)
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
			log.Println("No segment")
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
			log.Println(err)
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
			log.Println(err)
			return
		}

		_, err = tx.ExecContext(ctx, `INSERT INTO "segmentHistory" ("userID", "segmentName", operation, "dateReq") VALUES ($1, $2, $3, $4)`,
			request.UserID, val, 1, time.Now().Format(time.DateOnly))
		if err != nil {
			log.Println(err, 1)
			tx.Rollback()
			w.Write([]byte("Error: error in SQL request.\n"))
			return
		}

		_, err = tx.ExecContext(ctx, `INSERT INTO "userSegments"("userID", "segmentName") VALUES ($1, $2)`, request.UserID, val)
		if err != nil {
			log.Println(err, 2)
			tx.Rollback()
			w.Write([]byte("Error: error in SQL request.\n"))
			return
		}

		err = tx.Commit()
		if err != nil {
			log.Println(err, 3)
			w.Write([]byte("Error: error in SQL request.\n"))
			return
		}

		// ------
	}

	w.Write([]byte("Segments added.\n"))

	// deleting user from segments
	for _, val := range request.SegmetsToDelete {
		var n int
		// checking if user already in segment
		rows, err := DB.Query(`SELECT COUNT(*) FROM "userSegments" WHERE "segmentName" = $1`, val)
		if rows.Scan(&n); n == 0 || err != nil {
			continue
		}

		// SQL transaction to insert data in segmentHistory and userSegment
		// ------
		ctx := context.Background()
		tx, err := DB.BeginTx(ctx, nil)
		if err != nil {
			log.Println(err)
			return
		}

		_, err = tx.ExecContext(ctx, `INSERT INTO "segmentHistory" ("userID", "segmentName", operation, "dateReq") VALUES ($1, $2, $3, $4)`,
			request.UserID, val, 0, time.Now().Format(time.DateOnly))
		if err != nil {
			log.Println(err, 4)
			tx.Rollback()
			w.Write([]byte("Error: error in SQL request.\n"))
			return
		}

		_, err = tx.ExecContext(ctx, `DELETE FROM "userSegments" WHERE "segmentName" = $1`, val)
		if err != nil {
			log.Println(err, 5)
			tx.Rollback()
			w.Write([]byte("Error: error in SQL request.\n"))
			return
		}

		err = tx.Commit()
		if err != nil {
			log.Println(err)
			w.Write([]byte("Error: error in SQL request.\n"))
			return
		}

		// ------
	}

	w.Write([]byte("Segments deleted.\n"))
}

// searching for users active segments
func activeUserSegments(w http.ResponseWriter, r *http.Request) {
	var userID userRequest
	err := decode(&userID, w, r)
	if err != nil {
		return
	}
	rows, err := DB.Query(`SELECT "segmentName" FROM "userSegments" WHERE "userID" = $1`, userID.UserID)
	if err != nil {
		w.Write([]byte("Error: error in SQL request.\n"))
		log.Println(err)
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
		return
	}
	w.Write(response)
}

func userHistory(w http.ResponseWriter, r *http.Request) {
	var userReq userRequest
	err := decode(&userReq, w, r)
	if err != nil {
		return
	}

}
