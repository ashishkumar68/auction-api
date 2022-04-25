package forms

import (
	"github.com/ashishkumar68/auction-api/models"
	"time"
)

type EditItemForm struct {
	AuditableForm

	Item        *models.Item
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Category    models.ItemCategory `json:"category" binding:"oneof=0 1 2 3"`
	BrandName   string              `json:"brandName"`
	MarketValue models.Value        `json:"marketValue"`
	LastBidDate time.Time           `json:"lastBidDate" time_format:"2006-01-02"`
}
