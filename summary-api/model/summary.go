package model

import "time"

type Summary struct {
	ID        string    `json:"_id"`
	DialogID  string    `json:"dialogID"`
	Timestamp time.Time `json:"timestamp"`
	Summary   string    `json:"summary"`
}
