package client

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func MakeRequest(
	endpoint string,
	method string,
	addHead map[string]string,
	setHead map[string]string,
	timeout time.Duration,
	payload []byte,
) (*http.Response, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest(method, endpoint, bytes.NewReader(payload))
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
