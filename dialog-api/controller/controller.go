package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"main/client"
	"main/db"
	"main/model"
	"main/util"
	"net/http"
	"time"
)

type DialogController struct {
	db     db.CouchDB
	client client.Client
}

func NewDialogController(db db.CouchDB, client client.Client) DialogController {
	return DialogController{db: db, client: client}
}

func (receiver DialogController) UserID(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"userID": util.UserID})
}

func (receiver DialogController) DialogID(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"dialogID": uuid.New().String()})
}

func (receiver DialogController) AddDialog(ctx *gin.Context) {
	dialog, err := util.GetDialog(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.Error{Error: err.Error()})
		return
	}

	err = receiver.db.AddDialog(dialog)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}

	answer, err := receiver.client.GetAnswer(dialog)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}

	dialog.ID = uuid.New().String()
	dialog.DialogType = "bot"
	dialog.Timestamp = time.Now()
	dialog.Text = answer

	err = receiver.db.AddDialog(dialog)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"answer": answer})
}

func (receiver DialogController) EndDialog(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "dialog ended"})
}
