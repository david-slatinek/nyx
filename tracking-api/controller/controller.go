package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"main/db"
	"main/model"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type RecommendFollow struct {
	db db.RecommendFollow
}

func NewRecommendFollow(db db.RecommendFollow) RecommendFollow {
	return RecommendFollow{db: db}
}

func (receiver RecommendFollow) AddRecommendFollow(ctx *gin.Context) {
	fk, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.Error{Error: err.Error()})
		return
	}

	if fk <= 0 {
		ctx.JSON(http.StatusBadRequest, model.Error{Error: "fk must be greater than 0"})
		return
	}

	count := model.RecommendFollow{
		FkRecommend: fk,
		ClickAt:     time.Now(),
	}

	err = receiver.db.Create(count)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}

	if rand.Intn(11) <= 2 {
		order := model.Order{
			FkRecommend: fk,
			OrderAt:     time.Now(),
			Quantity:    rand.Intn(5-1) + 1,
		}
		err := receiver.db.CreateOrder(order)
		if err != nil {
			log.Printf("failed to create order: %v", err)
		}
	}

	ctx.HTML(http.StatusCreated, "index.html", gin.H{})
}
