package validators

import "github.com/go-playground/validator/v10"

const (
	TagImageFiles         = "images"
	TagMaxImageUploadSize = "max-img-file-size"
)

var customValidators = map[string]validator.Func{
	TagImageFiles:         ImageFiles,
	TagMaxImageUploadSize: MaxUploadImageFileSize,
}
