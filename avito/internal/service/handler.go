package service

import (
	//"encoding/json"
	//"log"
	"net/http"

	"github.com/Lalipopp4/avito/internal/db"
	"github.com/Lalipopp4/avito/internal/logger"
	"github.com/Lalipopp4/avito/internal/middleware"
	"github.com/gorilla/mux"
)

type segment struct {
	Name string `json:"segmentName"`
	Perc int    `json:"percent"`
}

type addUserRequest struct {
	SegmentsToAdd   []string `json:"segmentsToAdd"`
	SegmetsToDelete []string `json:"segmentsToDelete"`
	UserID          int      `json:"userID"`
	TTL             string   `json:"ttl"`
}

type userRequest struct {
	UserID int `json:"userID"`
}

type historyRequest struct {
	Date string `json:"date"`
}

// server handlers
func Handle() {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/addsegment", middleware.LoggingLay(addSegment, "Adding segment.")).Methods("POST")

	rtr.HandleFunc("/deletesegment", middleware.LoggingLay(deleteSegment, "Deleting segment.")).Methods("POST")

	rtr.HandleFunc("/adduser", middleware.LoggingLay(addUser, "Adding user, editting his segments.")).Methods("POST")

	rtr.HandleFunc("/activeusersegments", middleware.LoggingLay(activeUserSegments, "List of active user segments.")).Methods("POST")

	//rtr.HandleFunc("/userhistory", mw.MiddlewareLayer()).Methods("GET")

	go checkTTL()

	http.Handle("/", rtr)

}

// closing all connections
func Stop() {
	db.DB.Close()
	logger.Logger.Close()
}
