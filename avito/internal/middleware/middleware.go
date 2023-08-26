package middleware

import (
	// "log"
	//"encoding/json"
	"net/http"
	//"github.com/Lalipopp4/avito/internal/service"
)

type MiddlewareLayer struct{}

func SegmentCheck(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//
		next(w, r)
	}
}
