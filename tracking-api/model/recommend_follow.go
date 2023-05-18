package model

import "time"

type RecommendFollow struct {
	ID          int64     `gorm:"primaryKey"`
	FkRecommend int       `gorm:"column:fk_recommend"`
	ClickAt     time.Time `gorm:"column:click_at"`
}
