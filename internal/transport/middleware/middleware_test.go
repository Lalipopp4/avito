package middleware

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/Lalipopp4/test_api/pkg/logging"
	"github.com/gin-gonic/gin"
)

type pair struct {
	json string
	res  string
}

func TestValidateSegment(t *testing.T) {
	testReqs := []pair{{`{"segmentName": "E"}`, ""}, {`{"segmentName": "w", "Perc: 10}`, "Error in request.\n"}}
	for _, req := range testReqs {
		body, err := json.Marshal([]byte(req.json))
		if err != nil {
			t.Errorf("Error: json marshalling error")
		}
		rtr := gin.Default()
		logger, err := logging.New()
		if err != nil {
			t.Errorf("Error: New error")
		}
		var data []byte
		rtr.POST("segment/add", func(c *gin.Context) {
			c.Request.Body.Read(body)
			ValidateSegment(logger)
			data, err = io.ReadAll(c.Request.Response.Body)
			if err != nil {
				t.Errorf("Error: New error")
			}
		})

		if err != nil {
			t.Errorf("Expected error to be nil got %v", err)
		}
		if string(data) != req.res {
			t.Errorf("Expected "+req.res+" got %v", string(data))
		}
	}
}
