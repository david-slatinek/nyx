package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func GetDialogs(dialogID string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, os.Getenv("DIALOG_URL")+"/dialog/"+dialogID, nil)
	if err != nil {
		log.Printf("error creating http request: %s", err)
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error making http request: %s", err)
		return "", err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("could not read response body: %s", err)
		return "", err
	}

	var dialogs []map[string]string
	if err := json.Unmarshal(resBody, &dialogs); err != nil {
		log.Printf("failed to unmarshal message: %v", err)
		return "", err
	}

	var dialogText string
	for _, dialog := range dialogs {
		dialogText += fmt.Sprintf("%s: %s\n", dialog["dialogType"], dialog["text"])
	}
	return dialogText, nil
}
