package client

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func MakeRequest(
	endpoint string,
	method string,
	addHead map[string]string,
	setHead map[string]string,
	timeout time.Duration,
	payload io.Reader,
) (*http.Response, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest(method, endpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("got error %s", err.Error())
	}
	for key, val := range addHead {
		req.Header.Add(key, val)
	}
	for key, val := range setHead {
		req.Header.Set(key, val)
	}

	return client.Do(req)
}

func MakeMultiPartWriterFromFiles(fieldName string, files ...*os.File) (io.ReadWriter, string, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	for _, file := range files {
		image1Bytes, err := io.ReadAll(file)
		if err != nil {
			return nil, "", err
		}
		fi, err := file.Stat()
		if err != nil {
			return nil, "", err
		}
		err = file.Close()
		if err != nil {
			return nil, "", err
		}
		part, err := writer.CreateFormFile(fieldName, fi.Name())
		if err != nil {
			return nil, "", err
		}
		_, err = part.Write(image1Bytes)
		if err != nil {
			return nil, "", err
		}
	}
	err := writer.Close()
	if err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
}
