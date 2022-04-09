package forms

import (
	"github.com/ashishkumar68/auction-api/models"
)

type PlaceNewItemBidForm struct {

	ItemId		uint			`json:"itemId" binding:"required,min=1,max=12"`
	BidValue	models.Value	`json:"bidValue" binding:"required,min=1"`
	BidUser		*models.User	`json:"-" binding:"required"`
}
