package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Simple API Tester")
	fmt.Println("------------------")

	// 1. Get Method
	fmt.Print("Enter HTTP method (GET/POST): ")
	method, _ := reader.ReadString('\n')
	method = strings.TrimSpace(strings.ToUpper(method))
	if method != "GET" && method != "POST" {
		fmt.Println("Invalid method.  Using GET.")
		method = "GET"
	}

	// 2. Get URL
	fmt.Print("Enter URL: ")
	url, _ := reader.ReadString('\n')
	url = strings.TrimSpace(url)
	if url == "" {
		fmt.Println("URL cannot be empty.")
		return
	}

	// 3. Get Headers
	headers := make(map[string]string)
	for {
		fmt.Print("Enter header (key:value, or leave empty to finish): ")
		headerLine, _ := reader.ReadString('\n')
		headerLine = strings.TrimSpace(headerLine)

		if headerLine == "" {
			break
		}

		parts := strings.SplitN(headerLine, ":", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid header format.  Use key:value.")
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		headers[key] = value
	}

	// 4. Get Request Body
	fmt.Print("Enter request body (or leave empty for GET): ")
	body, _ := reader.ReadString('\n')
	body = strings.TrimSpace(body)

	// 5. Determine Content-Type (Basic)
	contentType := ""
	if method == "POST" && body != "" {
		contentType = "application/json"  // Default to JSON for POST requests
		if _, ok := headers["Content-Type"]; ok {
			contentType = headers["Content-Type"] // Override if set explicitly
		}
	}

	// 6. Prepare Request Parameters
	params := APIRequestParams{
		Method:      method,
		URL:         url,
		Headers:     headers,
		Body:        body,
		ContentType: contentType,
		Timeout:     10 * time.Second,
	}

	// 7. Make the API Request
	resp := MakeAPIRequest(params)

	// 8. Display the Response
	DisplayResponse(resp)
}