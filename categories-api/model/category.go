package model

type Category struct {
	ID            string     `json:"_id"`
	Name          string     `json:"name" binding:"required"`
	Subcategories []Category `json:"subcategories,omitempty"`
}
