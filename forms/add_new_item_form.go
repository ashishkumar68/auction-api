package forms

import (
	"github.com/ashishkumar68/auction-api/models"
	"time"
)

type AddNewItemForm struct {
	AuditableForm

	Name        string              `json:"name" binding:"required"`
	Description string              `json:"description" binding:"required"`
	Category    models.ItemCategory `json:"category" binding:"oneof=0 1 2 3"`
	BrandName   string              `json:"brandName" binding:"required"`
	MarketValue models.Value        `json:"marketValue" binding:"required"`
	LastBidDate time.Time           `json:"lastBidDate" binding:"required" time_format:"2006-01-02"`
}
