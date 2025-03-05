package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// APIResponse holds the response data from an API request.
type APIResponse struct {
	StatusCode int
	Headers    map[string][]string // Use map[string][]string for headers
	Body       string
	Error      error
	Duration   time.Duration
}

// MakeAPIRequest sends an API request and returns the response.
func MakeAPIRequest(method, apiURL string, headers map[string]string, body string) *APIResponse {
	start := time.Now()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	var bodyReader io.Reader
	if body != "" {
		bodyReader = bytes.NewBuffer([]byte(body))
	}

	req, err := http.NewRequest(method, apiURL, bodyReader)
	if err != nil {
		return &APIResponse{Error: fmt.Errorf("error creating request: %w", err)}
	}

	for key, value := range headers {
		req.Header.Set(key, value)
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

// parseHeaders parses the header string from the form into a map (OLD)
// Keeping it here for reference, but it's no longer used.
func parseHeaders(headerString string) (map[string]string, error) {
	headers := make(map[string]string)
	if headerString == "" {
		return headers, nil // Return an empty map if headerString is empty
	}

	headerPairs := strings.Split(headerString, "\n") // Split by newline character

	for _, pair := range headerPairs {
		pair = strings.TrimSpace(pair) // Trim leading/trailing spaces
		if pair == "" {
			continue // Skip empty lines
		}

		parts := strings.SplitN(pair, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header format: %s", pair)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		headers[key] = value
	}

	return headers, nil
}

func main() {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl.Execute(w, nil)
			return
		}

		if r.Method == http.MethodPost {
			// 1. Parse the form
			err := r.ParseForm()  // MUST call ParseForm before accessing form values with arrays
			if err != nil {
				http.Error(w, "Error parsing form", http.StatusBadRequest)
				return
			}

			// 2. Get form values
			method := r.FormValue("method")
			apiURL := r.FormValue("url")
			body := r.FormValue("body")

			// 3. Extract Headers from the form
			headerKeys := r.Form["header_key[]"]
			headerValues := r.Form["header_value[]"]

			// 4. Create Header Map
			headers := make(map[string]string)
			for i := 0; i < len(headerKeys); i++ {
				key := strings.TrimSpace(headerKeys[i])
				value := strings.TrimSpace(headerValues[i])
				if key != "" && value != "" { // Important: Check for empty keys/values
					headers[key] = value
				}
			}

			// Make API Request
			resp := MakeAPIRequest(method, apiURL, headers, body)

			// Prepare data for template
			data := struct {
				Method     string
				URL        string
				Headers    map[string][]string
				Body       string
				StatusCode int
				Error      string
				Duration   time.Duration
			}{
				Method:     method,
				URL:        apiURL,
				Headers:    resp.Headers,
				Body:       resp.Body,
				StatusCode: resp.StatusCode,
				Duration:   resp.Duration,
			}

			if resp.Error != nil {
				data.Error = resp.Error.Error()
			}

			// Execute template with API response
			tmpl.Execute(w, data)
		}
	})

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}