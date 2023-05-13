package controller

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"main/client"
	"main/db"
	"main/model"
	"main/util"
	"net/http"
	"sort"
	"time"
)

type DialogController struct {
	db       db.CouchDB
	client   client.Client
	svc      *sqs.SQS
	queueURL string
}

func NewDialogController(db db.CouchDB, client client.Client, svc *sqs.SQS) (DialogController, error) {
	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(util.QueueName),
	})
	if err != nil {
		return DialogController{}, err
	}

	return DialogController{
		db:       db,
		client:   client,
		svc:      svc,
		queueURL: *result.QueueUrl,
	}, nil
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
	dialog, err := util.GetDialog(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.Error{Error: err.Error()})
		return
	}

	jsonString, err := json.Marshal(gin.H{"dialogID": dialog.DialogID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}

	_, err = receiver.svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(jsonString)),
		QueueUrl:    aws.String(receiver.queueURL),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "dialog ended"})
}

func (receiver DialogController) GetDialog(ctx *gin.Context) {
	dialogID := ctx.Param("id")

	if len(dialogID) != 36 {
		ctx.JSON(http.StatusBadRequest, model.Error{Error: "invalid dialog id, must be 36 characters long"})
		return
	}

	dialogs, err := receiver.db.GetByDialogID(dialogID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}
	sort.Slice(dialogs, func(i, j int) bool {
		return dialogs[i].Timestamp.Before(dialogs[j].Timestamp)
	})
	ctx.JSON(http.StatusOK, dialogs)
}

func (receiver DialogController) GetDialogs(ctx *gin.Context) {
	dialogs, err := receiver.db.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dialogs)
}
