package service

import (
	"net/http"

	"github.com/Lalipopp4/avito/internal/logger"
	"github.com/Lalipopp4/avito/internal/middleware"
	"github.com/gorilla/mux"
)

// channel that is connected with ttl goroutine
var c = make(chan int, 1)

// server handlers
func Handle() {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/addsegment", middleware.LoggingLay(addSegment, "Adding segment.")).Methods("POST")

	rtr.HandleFunc("/deletesegment", middleware.LoggingLay(deleteSegment, "Deleting segment.")).Methods("POST")

	rtr.HandleFunc("/adduser", middleware.LoggingLay(addUser, "Adding user, editting his segments.")).Methods("POST")

	rtr.HandleFunc("/activeusersegments", middleware.LoggingLay(activeUserSegments, "List of active user segments.")).Methods("POST")

	rtr.HandleFunc("/segmenthistory", middleware.LoggingLay(segmentHistory, "Watching history of segments.")).Methods("POST")

	// url with csv file
	rtr.HandleFunc("/csvhistory", middleware.LoggingLay(csvHistory, "Watching csv file.")).Methods("GET")

	// goroutine that handles list of ttl users
	go checkTTL(c)

	http.Handle("/", rtr)

}

// closing all connections
func Stop() {
	c <- 1
	logger.Logger.Close()
}
