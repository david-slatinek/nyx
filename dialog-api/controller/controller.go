package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"main/db"
	"main/model"
	"main/util"
	"net/http"
	"time"
)

type DialogController struct {
	db db.CouchDB
}

func NewDialogController(db db.CouchDB) DialogController {
	return DialogController{db: db}
}

func (receiver DialogController) UserID(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"userID": util.UserID})
}

func (receiver DialogController) DialogID(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"dialogID": uuid.New().String()})
}

func (receiver DialogController) AddDialog(ctx *gin.Context) {
	var dialog model.Dialog
	if err := ctx.ShouldBindJSON(&dialog); err != nil {
		ctx.JSON(http.StatusBadRequest, model.Error{Error: err.Error()})
		return
	}

	dialog.ID = uuid.New().String()
	dialog.UserID = util.UserID

	if len(dialog.DialogID) != 36 {
		ctx.JSON(http.StatusBadRequest, model.Error{Error: "invalid dialog id, must be uuid v4"})
		return
	}
	dialog.DialogType = "user"
	dialog.Timestamp = time.Now()

	err := receiver.db.AddDialog(dialog)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}
	ctx.Status(http.StatusCreated)
}
