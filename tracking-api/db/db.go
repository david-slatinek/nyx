package db

import (
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"main/model"
	"os"
	"time"
)

type RecommendFollow struct {
	database *gorm.DB
}

func NewRecommendFollow() (RecommendFollow, error) {
	db, err := gorm.Open(mysql.Open(os.Getenv("MYSQL_URL")), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return RecommendFollow{}, err
	}
	return RecommendFollow{database: db}, err
}

func (receiver RecommendFollow) Close() error {
	db, err := receiver.database.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (receiver RecommendFollow) Create(count model.RecommendFollow) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return receiver.database.WithContext(ctx).Create(&count).Error
}
