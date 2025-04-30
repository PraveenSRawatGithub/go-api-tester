package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl.Execute(w, nil)
			return
		}

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Error parsing form", http.StatusBadRequest)
				return
			}

			method := r.FormValue("method")
			apiURL := r.FormValue("url")
			body := r.FormValue("body")

			// Extract Headers from form
			headerKeys := r.Form["header_key[]"]
			headerValues := r.Form["header_value[]"]
			headers := make(map[string]string)
			for i := 0; i < len(headerKeys); i++ {
				key := strings.TrimSpace(headerKeys[i])
				value := strings.TrimSpace(headerValues[i])
				if key != "" && value != "" {
					headers[key] = value
				}
			}

			// Use the updated MakeAPIRequest from api.go
			params := APIRequestParams{
				Method:      method,
				URL:         apiURL,
				Headers:     headers,
				Body:        body,
				ContentType: headers["Content-Type"],
				Timeout:     10 * time.Second,
			}

			resp := MakeAPIRequest(params)

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

			tmpl.Execute(w, data)
		}
	})

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
