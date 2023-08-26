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

func Handle() {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/addsegment", middleware.SegmentCheck(addSegment)).Methods("POST")

	rtr.HandleFunc("/deletesegment", middleware.SegmentCheck(deleteSegment)).Methods("POST")

	rtr.HandleFunc("/adduser", addUser).Methods("POST")

	rtr.HandleFunc("/activeusersegments", activeUserSegments).Methods("POST")

	//rtr.HandleFunc("/userhistory", mw.MiddlewareLayer()).Methods("GET")

	http.Handle("/", rtr)

}

func Stop() {
	db.Close(DB)
	logger.Logger.Close()
}
