package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type pair struct {
	json string
	res  string
}

func TestAddSegment(t *testing.T) {
	testReqs := []pair{{`{"segmentName": "E"}`, "E added."}, {`{"segmentName": "w", "Perc: 10}`, "Error: this segment is already exists or empty request.\n"}}
	for _, req := range testReqs {
		body, err := json.Marshal([]byte(req.json))
		if err != nil {
			t.Errorf("Error: json marshalling error")
		}
		r := httptest.NewRequest(http.MethodPost, "/addsegment", bytes.NewReader(body))
		defer r.Body.Close()
		w := httptest.NewRecorder()
		addSegment(w, r)
		res := w.Result()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Expected error to be nil got %v", err)
		}
		if string(data) != req.res {
			t.Errorf("Expected "+req.res+" got %v", string(data))
		}
	}
}

func TestDeleteSegment(t *testing.T) {
	testReqs := []pair{{`{"segmentName": "b"}`, "b deleted."}, {`{"segmentName": "w"}`, "Error: this segment doesn't exist or empty request.\n"}}
	for _, req := range testReqs {
		body, err := json.Marshal([]byte(req.json))
		if err != nil {
			t.Errorf("Error: json marshalling error")
		}
		r := httptest.NewRequest(http.MethodPost, "/deletesegment", bytes.NewReader(body))
		defer r.Body.Close()
		w := httptest.NewRecorder()
		addSegment(w, r)
		res := w.Result()
		data, err := io.ReadAll(res.Body)
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
		{`{"user_id": 1, "segmentsToAdd": ["E", "w], "ttl": "2", "segmentsToDelete": ["D"]}`, "E added."},
		{`{"user_id" : 10, "segmentName": "w", "Perc: 10}`, "Error: this segment is already exists or empty request.\n"}}
	for _, req := range testReqs {
		body, err := json.Marshal([]byte(req.json))
		if err != nil {
			t.Errorf("Error: json marshalling error")
		}
		r := httptest.NewRequest(http.MethodPost, "/adduser", bytes.NewReader(body))
		defer r.Body.Close()
		w := httptest.NewRecorder()
		addSegment(w, r)
		res := w.Result()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Expected error to be nil got %v", err)
		}
		if string(data) != req.res {
			t.Errorf("Expected "+req.res+" got %v", string(data))
		}
	}
}

func TestActiveUserSegments(t *testing.T) {
	testReqs := []pair{
		{`{"user_id": 1`, `["E", "W"].`},
		{`{"user_id": 22}`, "[]"}}
	for _, req := range testReqs {
		body, err := json.Marshal([]byte(req.json))
		if err != nil {
			t.Errorf("Error: json marshalling error")
		}
		r := httptest.NewRequest(http.MethodPost, "/activeuserhistory", bytes.NewReader(body))
		defer r.Body.Close()
		w := httptest.NewRecorder()
		addSegment(w, r)
		res := w.Result()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Expected error to be nil got %v", err)
		}
		if string(data) != req.res {
			t.Errorf("Expected "+req.res+" got %v", string(data))
		}
	}
}

func TestSegmentHistory(t *testing.T) {
	testReqs := []pair{{`"date": "2022-08`, `No records on "2022-08".`}, {`"date": "2023-08`, `/csvhistory?date=2023-08".`}}
	for _, req := range testReqs {
		body, err := json.Marshal([]byte(req.json))
		if err != nil {
			t.Errorf("Error: json marshalling error")
		}
		r := httptest.NewRequest(http.MethodPost, "/segmenthistory", bytes.NewReader(body))
		defer r.Body.Close()
		w := httptest.NewRecorder()
		addSegment(w, r)
		res := w.Result()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Expected error to be nil got %v", err)
		}
		if string(data) != req.res {
			t.Errorf("Expected "+req.res+" got %v", string(data))
		}
	}
}
