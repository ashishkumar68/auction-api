package forms

import (
	"github.com/ashishkumar68/auction-api/models"
)

type PlaceNewItemBidForm struct {
	AuditableForm

	ItemId    uint         `json:"itemId" binding:"required,min=1,max=12"`
	BidValue  models.Value `json:"bidValue" binding:"required,min=1"`
	BidUserId uint
}
