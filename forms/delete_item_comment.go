package forms

import "github.com/ashishkumar68/auction-api/models"

type DeleteItemCommentForm struct {
	AuditableForm

	Item	*models.Item
	Comment	*models.ItemComment
}
