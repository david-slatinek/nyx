package model

type Recommendation struct {
	ID    string  `json:"id"`
	Label string  `json:"label"`
	Score float32 `json:"score"`
}
