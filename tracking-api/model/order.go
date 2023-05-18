package model

import "time"

type Order struct {
	ID          int       `gorm:"primaryKey"`
	FkRecommend int       `gorm:"column:fk_recommend"`
	OrderAt     time.Time `gorm:"column:order_at"`
	Quantity    int       `gorm:"column:quantity"`
}

func (Order) TableName() string {
	return "order_table"
}
