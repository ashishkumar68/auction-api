package validators

import (
	"github.com/ashishkumar68/auction-api/utils"
	"github.com/go-playground/validator/v10"
	"math"
	"mime/multipart"
	"strconv"
)

var allowedImageMimeTypes = []string{
	"image/jpeg", "image/png",
}

var ImageFiles validator.Func = func(fl validator.FieldLevel) bool {
	files, ok := fl.Field().Interface().([]*multipart.FileHeader)
	if !ok {
		return false
	}

	for _, file := range files {
		contentType := file.Header.Get("Content-Type")
		for _, allowedType := range allowedImageMimeTypes {
			if allowedType == contentType {
				return true
			}
		}

		return false
	}

	return true
}

var MaxUploadImageFileSize validator.Func = func(fl validator.FieldLevel) bool {
	files, ok := fl.Field().Interface().([]*multipart.FileHeader)
	if !ok {
		return false
	}

	maxAllowedSize, err := strconv.Atoi(fl.Param())
	utils.PanicIf(err)

	for _, file := range files {
		// converting Bytes to MBs
		fileSizeMBs := math.Round(float64(file.Size) / 1024 / 1024)
		if fileSizeMBs > float64(maxAllowedSize) {
			return false
		}
	}

	return true
}
