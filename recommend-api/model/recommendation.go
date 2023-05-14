package model

type Recommendation struct {
	ID     string  `json:"id"`
	Dialog string  `json:"dialog"`
	Label  string  `json:"labels"`
	Score  float32 `json:"scores"`
}
