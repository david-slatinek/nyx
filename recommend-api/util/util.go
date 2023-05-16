package util

import (
	"github.com/levigross/grequests"
	"log"
	"main/model"
	"os"
	"time"
)

var Categories []model.Category
var CategoriesText []string

func GetDialogs(dialogID string) ([]model.Dialog, error) {
	resp, err := grequests.Get(os.Getenv("DIALOG_URL")+"/dialog/"+dialogID, nil)
	if err != nil {
		return nil, err
	}

	var dialogs []model.Dialog
	err = resp.JSON(&dialogs)
	return dialogs, err
}

func GetCategories() {
	log.Printf("Getting categories at %v", time.Now().Format("2006-01-02 15:04:05"))

	resp, err := grequests.Get(os.Getenv("CATEGORY_URL")+"/categories", nil)
	if err != nil {
		log.Printf("failed to get categories: %v", err)
		return
	}

	var categories []model.Category
	err = resp.JSON(&categories)
	if err != nil {
		log.Printf("failed to get unmarshall categories: %v", err)
		return
	}
	Categories = categories

	getCategoriesNames(categories, &CategoriesText)
}

func getCategoriesNames(category []model.Category, categories *[]string) {
	for i := range category {
		*categories = append(*categories, category[i].Name)

		if category[i].Subcategories != nil {
			getCategoriesNames(category[i].Subcategories, categories)
		}
	}
}

func GetMainCategoryID(categories []model.Category, label string) string {
	for _, category := range categories {
		if category.Name == label {
			return category.ID
		}
	}

	for _, category := range categories {
		if category.Name == label {
			return category.ID
		}

		if category.Subcategories != nil {
			return GetMainCategoryID(category.Subcategories, label)
		}
	}

	return ""
}
