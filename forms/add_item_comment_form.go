package forms

import "github.com/ashishkumar68/auction-api/models"

type AddItemCommentForm struct {
	AuditableForm

	Item	*models.Item
	Comment	string	`json:"comment" binding:"required,min=3"`
}
