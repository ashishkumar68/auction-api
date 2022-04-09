package forms

import (
	"github.com/ashishkumar68/auction-api/models"
)

type AddNewItemForm struct {

	Name        string              `json:"name" binding:"required"`
	Description string              `json:"description" binding:"required"`
	Category    models.ItemCategory `json:"category" binding:"required,oneof=0 1 2 3"`
	BrandName   string              `json:"brandName" binding:"required"`
	MarketValue models.Value        `json:"marketValue" binding:"required"`
}