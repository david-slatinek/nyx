package controller

import (
	"github.com/gin-gonic/gin"
	"main/db"
	"main/model"
	"net/http"
	"time"
)

type RecommendFollow struct {
	db db.RecommendFollow
}

func NewRecommendFollow(db db.RecommendFollow) RecommendFollow {
	return RecommendFollow{db: db}
}

func (receiver RecommendFollow) AddRecommendFollow(ctx *gin.Context) {
	var fk model.FkRecommend
	if err := ctx.ShouldBindJSON(&fk); err != nil {
		ctx.JSON(http.StatusBadRequest, model.Error{Error: err.Error()})
		return
	}

	if fk.FkRecommend <= 0 {
		ctx.JSON(http.StatusBadRequest, model.Error{Error: "fk must be greater than 0"})
		return
	}

	count := model.RecommendFollow{
		FkRecommend: fk.FkRecommend,
		ClickAt:     time.Now(),
	}

	err := receiver.db.Create(count)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}
	ctx.HTML(http.StatusCreated, "index.html", gin.H{})
}
