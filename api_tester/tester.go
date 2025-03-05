package api_tester

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// SendRequest sends an API request (GET, POST, PUT, DELETE)
func SendRequest(method, url string, headers map[string]string, body []byte) ([]byte, error) {
	client := &http.Client{}

	// Create a new request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Add headers if provided
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
