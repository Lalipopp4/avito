package models

type User struct {
	Id int `json:"id"`
}

type UserRequest struct {
	Id               int      `json:"id"`
	SegmentsToAdd    []string `json:"segmentsToAdd"`
	SegmentsToDelete []string `json:"segmentsToDelete"`
	TTL              string   `json:"ttl"`
}
