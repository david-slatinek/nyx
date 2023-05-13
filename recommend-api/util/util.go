package util

import (
	"github.com/levigross/grequests"
	"main/model"
	"os"
)

func GetDialogs(dialogID string) ([]model.Dialog, error) {
	resp, err := grequests.Get(os.Getenv("DIALOG_URL")+"/dialog/"+dialogID, nil)
	if err != nil {
		return nil, err
	}

	var dialogs []model.Dialog
	err = resp.JSON(&dialogs)
	return dialogs, err
}
