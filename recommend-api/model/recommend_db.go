package model

import "time"

type RecommendDB struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	FkMainCategory string    `json:"fk_main_category"`
	FkSubCategory  string    `json:"fk_sub_category"`
	CategoryName   string    `json:"category_name"`
	Score          float32   `json:"score"`
	FkMainDialog   string    `json:"fk_main_dialog"`
	FkDialog       string    `json:"fk_dialog"`
	RecommendedAt  time.Time `json:"recommended_at"`
}
