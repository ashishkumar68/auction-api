package validators

import (
	"github.com/ashishkumar68/auction-api/utils"
	"github.com/go-playground/validator/v10"
	"math"
	"mime/multipart"
	"strconv"
	"strings"
)

var allowedImageMimeTypes = []string{
	"image/jpeg", "image/png",
}
var allowedExtensions = []string{
	"jpeg", "jpg", "png",
}

var ImageFiles validator.Func = func(fl validator.FieldLevel) bool {
	files, ok := fl.Field().Interface().([]*multipart.FileHeader)
	if !ok {
		return false
	}

	validFilesCount := 0
IterateFiles:
	for _, file := range files {
		contentType := file.Header.Get("Content-Type")
		fileNameInfo := strings.Split(file.Filename, ".")
		if len(fileNameInfo) < 2 {
			continue
		}
		extension := fileNameInfo[len(fileNameInfo)-1]
		for _, allowedType := range allowedImageMimeTypes {
			if allowedType == contentType {
				validFilesCount += 1
				continue IterateFiles
			}
		}
		for _, ext := range allowedExtensions {
			if ext == extension {
				validFilesCount += 1
				continue IterateFiles
			}
		}
	}

	if validFilesCount == len(files) {
		return true
	} else {
		return false
	}
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

var MaxUploadFilesCount validator.Func = func(fl validator.FieldLevel) bool {
	files, ok := fl.Field().Interface().([]*multipart.FileHeader)
	if !ok {
		return false
	}

	maxAllowedCount, err := strconv.Atoi(fl.Param())
	utils.PanicIf(err)

	if maxAllowedCount < len(files) {
		return false
	}

	return true
}
