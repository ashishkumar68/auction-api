package forms

import (
	"github.com/ashishkumar68/auction-api/models"
	"mime/multipart"
)

type AddItemImagesForm struct {
	AuditableForm

	Item       *models.Item            `form:"-"`
	ImageFiles []*multipart.FileHeader `form:"images" binding:"images,max-img-file-size=12"`
}
