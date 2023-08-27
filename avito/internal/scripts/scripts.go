package scripts

import (
	"encoding/json"
	"net/http"

	"github.com/Lalipopp4/avito/internal/logger"
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
