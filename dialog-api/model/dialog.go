package model

import "time"

type Dialog struct {
	ID         string    `json:"_id"`
	UserID     string    `json:"userID"`
	DialogID   string    `json:"dialogID" binding:"required"`
	DialogType string    `json:"dialogType"`
	Timestamp  time.Time `json:"timestamp"`
	Text       string    `json:"text" binding:"required"`
}
