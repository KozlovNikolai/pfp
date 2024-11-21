package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

func DoRequest[R any](method, url string, body io.Reader, headers ...map[string]string) (resp R, err error) {
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return
	}

	for _, head := range headers {
		for k, v := range head {
			req.Header.Set(k, v)
		}
	}

	client := http.Client{Timeout: time.Minute * 5}
	response, err := client.Do(req)

	if err != nil {
		return
	}

	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
	} else {
		err = decoder.Decode(&resp)
	}

	return
}
