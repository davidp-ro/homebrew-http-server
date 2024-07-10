// Description: A simple HTTP client that sends a POST request with some basic
// data to test the homebrew-http-server.

package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

const URL = "http://localhost:8080"

func post(contentType string, data []byte) {
	response, err := http.Post(URL, contentType, bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Request failed: %s\n", err)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %s\n", err)
		return
	}

	fmt.Printf("Response: %s\n", body)
}

func main() {
	// JSON body:
	post("application/json", []byte(`{"key": "value"}`))

	// Plain text body:
	post("plain/text", []byte(`Hello world!\nOn multiple lines\n\n\n\n\nAnd even more!`))
}
