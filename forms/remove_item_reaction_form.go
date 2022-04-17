package forms

import "github.com/ashishkumar68/auction-api/models"

type RemoveItemReactionForm struct {
	AuditableForm

	Item *models.Item `binding:"required"`
}
