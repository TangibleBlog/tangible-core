package net

import (
	"io"
	"net/http"
	"strings"
)

func Post(url string, contentType string, body []byte) ([]byte, error) {
	res, err := http.Post(url, contentType, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}
