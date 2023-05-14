package model

type RecommendResult struct {
	ID     string    `json:"id"`
	Dialog string    `json:"dialog"`
	Labels []string  `json:"labels"`
	Scores []float32 `json:"scores"`
}
