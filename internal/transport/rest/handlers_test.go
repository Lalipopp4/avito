package rest

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/gin-gonic/gin"
)

type pair struct {
	json string
	res  string
}

func TestAddSegment(t *testing.T) {
	testReqs := []pair{{`{"segmentName": "E"}`, "E added successfully.."}, {`{"segmentName": "w", "Perc: 10}`, "Error: this segment is already exists or empty request.\n"}}
	for _, req := range testReqs {
		body, err := json.Marshal([]byte(req.json))
		if err != nil {
			t.Errorf("Error: json marshalling error")
		}
		rtr := gin.Default()
		rAPI, err := New()
		if err != nil {
			t.Errorf("Error: New error")
		}
		var data []byte
		rtr.POST("segment/add", func(c *gin.Context) {
			c.Request.Body.Read(body)
			rAPI.addSegment(c)
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

func TestDeleteSegment(t *testing.T) {
	testReqs := []pair{{`{"segmentName": "E"}`, "E deleted successfully."}, {`{"segmentName": "f"}`, "f doesn't exist.\n"}}
	for _, req := range testReqs {
		body, err := json.Marshal([]byte(req.json))
		if err != nil {
			t.Errorf("Error: json marshalling error")
		}
		rtr := gin.Default()
		rAPI, err := New()
		if err != nil {
			t.Errorf("Error: New error")
		}
		var data []byte
		rtr.POST("segment/delete", func(c *gin.Context) {
			c.Request.Body.Read(body)
			rAPI.deleteSegment(c)
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

func TestGetSegmentsByUser(t *testing.T) {
	testReqs := []pair{{`{"id: 1"}`, "AVITO_10 AVITO_NO"}, {`{"id":10}`, "\n"}}
	for _, req := range testReqs {
		body, err := json.Marshal([]byte(req.json))
		if err != nil {
			t.Errorf("Error: json marshalling error")
		}
		rtr := gin.Default()
		rAPI, err := New()
		if err != nil {
			t.Errorf("Error: New error")
		}
		var data []byte
		rtr.GET("user/segments", func(c *gin.Context) {
			c.Request.Body.Read(body)
			rAPI.getSegmentsByUser(c)
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

func TestAddUser(t *testing.T) {
	testReqs := []pair{
		{`{"user_id": 1, "segmentsToAdd": ["E", "w], "ttl": "2", "segmentsToDelete": ["D"]}`, "Segments were added and deleted for user successfully."},
		{`{"user_id" : 10, "segmentsToDelete": ["w"]}`, "Segemnt w doesn't exists.\n"}}
	for _, req := range testReqs {
		body, err := json.Marshal([]byte(req.json))
		if err != nil {
			t.Errorf("Error: json marshalling error")
		}
		rtr := gin.Default()
		rAPI, err := New()
		if err != nil {
			t.Errorf("Error: New error")
		}
		var data []byte
		rtr.POST("user/add", func(c *gin.Context) {
			c.Request.Body.Read(body)
			rAPI.addSegment(c)
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

func TestGetHistoryByDate(t *testing.T) {
	testReqs := []pair{{`"date": "2022-08`, `No records on "2022-08".`}, {`"date": "2023-08`, `/csvhistory?date=2023-08".`}}
	for _, req := range testReqs {
		body, err := json.Marshal([]byte(req.json))
		if err != nil {
			t.Errorf("Error: json marshalling error")
		}
		rtr := gin.Default()
		rAPI, err := New()
		if err != nil {
			t.Errorf("Error: New error")
		}
		var data []byte
		rtr.POST("segment/history", func(c *gin.Context) {
			c.Request.Body.Read(body)
			rAPI.addSegment(c)
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
