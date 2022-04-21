package forms

import "github.com/ashishkumar68/auction-api/models"

type RemoveItemReactionForm struct {
	AuditableForm

	Item *models.Item `json:"-"`
}

func NewRemoveItemReactionForm(actionUser *models.User, item *models.Item) *RemoveItemReactionForm {
	return &RemoveItemReactionForm{
		AuditableForm: AuditableForm{
			ActionUser: actionUser,
		},
		Item:          item,
	}
}