package models

type Segment struct {
	Id   int
	Name string `json:"segmentName"`
	Perc int    `json:"perc"`
}
