package middleware

import (
	"net/http"

	"github.com/Lalipopp4/avito/internal/logger"
)

// implementation of logging middleware layer
func LoggingLay(next http.HandlerFunc, step string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Logger.Log()
		next(w, r)
	}
}
