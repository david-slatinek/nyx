package main

import (
	"log"
	"main/env"
	"main/model"
	"main/util"
)

func main() {
	err := env.Load("env/.env")
	if err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	recommend := model.Recommend{
		DialogID: "86ce24c7-8108-4d2f-86fe-a687f246c0d6",
		Summary:  "Sarah is in her early twenties and doesn't know what to do with herself. She doesn't have any siblings. She is not sure if she has any plans for the future. She's not sure what her name is. She's called Sarah.",
	}

	dialogs, err := util.GetDialogs(recommend.DialogID)
	if err != nil {
		log.Fatalf("failed to get dialogs: %v", err)
	}
	log.Printf("dialogs: %v", dialogs)
}
