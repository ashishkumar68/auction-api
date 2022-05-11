package forms

import "github.com/ashishkumar68/auction-api/models"

type RemoveItemThumbnailForm struct {
	AuditableForm

	Item *models.Item `form:"-"`
}
