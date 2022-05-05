package forms

import "github.com/ashishkumar68/auction-api/models"

type RemoveItemImageForm struct {
	AuditableForm

	ItemImage *models.ItemImage `form:"-"`
}
