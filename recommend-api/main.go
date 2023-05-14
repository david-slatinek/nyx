package main

import (
	"log"
	"main/client"
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
		ID:       "c3c27bd5-22ee-4c25-9f32-b8f5d0bedf40",
		DialogID: "86ce24c7-8108-4d2f-86fe-a687f246c0d6",
		Summary:  "Sarah is in her early twenties and doesn't know what to do with herself. She doesn't have any siblings. She is not sure if she has any plans for the future. She's not sure what her name is. She's called Sarah.",
	}

	dialogs, err := util.GetDialogs(recommend.DialogID)
	if err != nil {
		log.Fatalf("failed to get dialogs: %v", err)
	}

	categories, err := util.GetCategories()
	if err != nil {
		log.Fatalf("failed to get categories: %v", err)
	}

	var categoriesText = make([]string, 0, len(categories))
	util.GetCategoriesNames(categories, &categoriesText)

	rClient, err := client.NewClient()
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer func(rClient client.Client) {
		if err := rClient.Close(); err != nil {
			log.Printf("failed to close client: %v", err)
		}
	}(rClient)

	recommendResult, err := rClient.GetRecommendationDialogs(dialogs, categoriesText)
	if err != nil {
		log.Fatalf("failed to get recommendation for dialogs: %v", err)
	}
	log.Printf("recommendResult: %v", recommendResult)

	recommendResultSummary, err := rClient.GetRecommendationSummary(recommend, categoriesText)
	if err != nil {
		log.Fatalf("failed to get recommendation for summary: %v", err)
	}
	log.Printf("recommendResult: %v", recommendResultSummary)
}