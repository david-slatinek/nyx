package util

import (
	"github.com/levigross/grequests"
	"log"
	"main/model"
	"os"
	"time"
)

var CategoriesText []string
var CategoriesMap = make(map[string]string)

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

	getCategoriesNames(categories, &CategoriesText)
}

func getCategoriesNames(category []model.Category, categories *[]string) {
	for key, value := range category {
		*categories = append(*categories, category[key].Name)
		CategoriesMap[value.Name] = value.ID

		if category[key].Subcategories != nil {
			getCategoriesNames(category[key].Subcategories, categories)
		}
	}
}
