package utils

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

func PanicIf(err error) {
	if err == nil {
		return
	}

	panic(err)
}

func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	err = CreateAllDirectoriesForPath(dst)
	if err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func CreateAllDirectoriesForPath(filePath string) error {
	filePathMeta := strings.Split(filePath, "/")
	directory := strings.Join(filePathMeta[:len(filePathMeta)-1], "/")

	err := os.MkdirAll(directory, 0755)
	return err
}

func GetRenamedFileName(fileName string) (string, error) {
	fileNameMeta := strings.Split(fileName, ".")
	fileNameInfoCount := len(fileNameMeta)
	if 1 == fileNameInfoCount {
		return "", fmt.Errorf("invalid file name received: %s", fileName)
	}
	secondLastPieceIndex := fileNameInfoCount - 2
	fileNameMeta[secondLastPieceIndex] = fmt.Sprintf("%s_%s", fileNameMeta[secondLastPieceIndex], uuid.NewString())

	return strings.Join(fileNameMeta, "."), nil
}

func GetFileSystemFilePath(filePath string) string {
	return fmt.Sprintf("%s/%s", os.Getenv("FILE_UPLOADS_DIR"), filePath)
}

func GetGlobalUploadsDir() string {
	return os.Getenv("FILE_UPLOADS_DIR")
}
