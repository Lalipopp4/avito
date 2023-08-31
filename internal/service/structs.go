package service

type segment struct {
	Name string `json:"segmentName"`
	Perc int    `json:"percent"`
}

type addUserRequest struct {
	SegmentsToAdd    []string `json:"segmentsToAdd"`
	SegmentsToDelete []string `json:"segmentsToDelete"`
	UserID           int      `json:"userID"`
	TTL              string   `json:"ttl"`
}

type userRequest struct {
	UserID int `json:"userID"`
}

type historyRequest struct {
	Date string `json:"date"`
}

type point struct {
	userID      int
	segmentName string
}
