package model

import "time"

type RecommendDB struct {
	ID            int       `json:"id"`
	UserID        string    `json:"user_id"`
	FkCategory    string    `json:"fk_category"`
	CategoryName  string    `json:"category_name"`
	Score         float32   `json:"score"`
	FkMainDialog  string    `json:"fk_main_dialog"`
	FkDialog      string    `json:"fk_dialog"`
	RecommendedAt time.Time `json:"recommended_at"`
}

func (RecommendDB) TableName() string {
	return "recommend"
}
