package forms

import "github.com/ashishkumar68/auction-api/models"

type MarkItemOffBidForm struct {
	AuditableForm

	Item *models.Item
}
