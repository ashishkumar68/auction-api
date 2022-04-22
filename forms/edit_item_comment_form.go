package forms

import "github.com/ashishkumar68/auction-api/models"

type EditItemCommentForm struct {
	AuditableForm

	Item			*models.Item
	Comment			*models.ItemComment
	EditedComment	string	`json:"comment" binding:"required,min=3"`
}
