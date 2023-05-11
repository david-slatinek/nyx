package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"main/db"
	"main/model"
	"net/http"
)

type CategoryController struct {
	db db.CouchDB
}

func NewCategoryController(db db.CouchDB) CategoryController {
	return CategoryController{
		db: db,
	}
}

func setIDs(category *model.Category) {
	category.ID = uuid.New().String()
	for i := range category.Subcategories {
		category.Subcategories[i].ID = uuid.New().String()

		if category.Subcategories[i].Subcategories != nil {
			setIDs(&category.Subcategories[i])
		}
	}
}

func (receiver CategoryController) AddCategory(ctx *gin.Context) {
	var category model.Category
	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, model.Error{Error: err.Error()})
		return
	}

	setIDs(&category)

	if err := receiver.db.AddCategory(category); err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, category)
}
