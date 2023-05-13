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

func GetCategories() ([]model.Category, error) {
	resp, err := grequests.Get(os.Getenv("CATEGORY_URL")+"/categories", nil)
	if err != nil {
		return nil, err
	}

	var categories []model.Category
	err = resp.JSON(&categories)
	return categories, err
}

func GetCategoriesNames(category []model.Category, categories *[]string) {
	for i := range category {
		*categories = append(*categories, category[i].Name)

		if category[i].Subcategories != nil {
			GetCategoriesNames(category[i].Subcategories, categories)
		}
	}
}
