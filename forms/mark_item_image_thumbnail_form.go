package forms

import "github.com/ashishkumar68/auction-api/models"

type MarkItemImageThumbnailForm struct {
	AuditableForm

	ItemImg *models.ItemImage `form:"-"`
}
