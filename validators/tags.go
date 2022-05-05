package validators

import "github.com/go-playground/validator/v10"

const (
	TagImageFiles            = "images"
	TagMaxImageUploadSize    = "max-img-file-size"
	TagMaxUploadAllowedFiles = "max-upload-files"
)

var customValidators = map[string]validator.Func{
	TagImageFiles:            ImageFiles,
	TagMaxImageUploadSize:    MaxUploadImageFileSize,
	TagMaxUploadAllowedFiles: MaxUploadFilesCount,
}
