package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// APIRequestParams holds the parameters for an API request.
type APIRequestParams struct {
	Method      string
	URL         string
	Headers     map[string]string
	Body        string
	ContentType string // Content-Type of the request body
	Timeout     time.Duration
}

// APIResponse holds the response data from an API request.
type APIResponse struct {
	StatusCode int
	Headers    map[string][]string // Use map[string][]string for headers
	Body       string
	Error      error
	Duration   time.Duration
}

// MakeAPIRequest sends an API request and returns the response.  It now takes APIRequestParams.
func MakeAPIRequest(params APIRequestParams) *APIResponse {
	start := time.Now()
	client := &http.Client{
		Timeout: params.Timeout,
	}

	body := bytes.NewBuffer([]byte(params.Body))

	req, err := http.NewRequest(params.Method, params.URL, body)
	if err != nil {
		return &APIResponse{Error: fmt.Errorf("error creating request: %w", err)}
	}

	for key, value := range params.Headers {
		req.Header.Set(key, value)
	}

	if params.ContentType != "" {
		req.Header.Set("Content-Type", params.ContentType)
	}

	resp, err := client.Do(req)
	if err != nil {
		return &APIResponse{Error: fmt.Errorf("error sending request: %w", err)}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &APIResponse{Error: fmt.Errorf("error reading response: %w", err)}
	}

	elapsed := time.Since(start)

	return &APIResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       string(respBody),
		Error:      nil,
		Duration:   elapsed,
	}
}

// FormatJSON formats a JSON string for readability.
func FormatJSON(jsonString string) (string, error) {
	var jsonData interface{}
	err := json.Unmarshal([]byte(jsonString), &jsonData)
	if err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}

	prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error formatting JSON: %w", err)
	}

	return string(prettyJSON), nil
}

// DisplayResponse prints the API response to the console.
func DisplayResponse(resp *APIResponse) {
	if resp.Error != nil {
		fmt.Println("Error:", resp.Error)
		return
	}

	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Response Time:", resp.Duration)
	fmt.Println("Headers:")
	for key, values := range resp.Headers { // Iterate over the slice of values
		fmt.Printf("  %s: %s\n", key, strings.Join(values, ", "))
	}

	fmt.Println("Body:")
	if strings.Contains(strings.Join(resp.Headers["Content-Type"],","), "application/json") {
		formattedJSON, err := FormatJSON(resp.Body)
		if err == nil {
			fmt.Println(formattedJSON)
		} else {
			fmt.Println(resp.Body)
		}
	} else {
		fmt.Println(resp.Body)
	}
}