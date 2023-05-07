package model

import "time"

type Chat struct {
	ID        string    `json:"_id"`
	UserID    string    `json:"userID" binding:"required"`
	ChatID    string    `json:"chatID" binding:"required"`
	ChatType  string    `json:"chatType" binding:"required"`
	Timestamp time.Time `json:"timestamp" time_format:"2006-01-02 15:04:05"`
	Text      string    `json:"text" binding:"required"`
}
