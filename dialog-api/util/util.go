package util

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"main/model"
	"time"
)

const UserID = "93cb8a64-97bc-4742-acba-a1fdbc543434"
const QueueName = "summary-q"

func GetDialog(ctx *gin.Context) (model.Dialog, error) {
	var dialog model.Dialog
	if err := ctx.ShouldBindJSON(&dialog); err != nil {
		return model.Dialog{}, err
	}

	dialog.ID = uuid.New().String()
	dialog.UserID = UserID

	if len(dialog.DialogID) != 36 {
		return model.Dialog{}, errors.New("invalid dialog id, must be 36 characters long")
	}
	dialog.DialogType = "user"
	dialog.Timestamp = time.Now()

	return dialog, nil
}
