package model

type RecommendResult struct {
	Dialog string    `json:"dialog"`
	Labels []string  `json:"labels"`
	Scores []float32 `json:"scores"`
}
