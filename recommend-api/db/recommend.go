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

type RecommendDB struct {
	database *gorm.DB
}

func NewRecommendDB() (RecommendDB, error) {
	db, err := gorm.Open(mysql.Open(os.Getenv("MYSQL_URL")), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return RecommendDB{}, err
	}
	return RecommendDB{database: db}, err
}

func (r RecommendDB) Close() error {
	db, err := r.database.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (r RecommendDB) Create(recommend []model.RecommendDB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.database.WithContext(ctx).Create(&recommend).Error
}
