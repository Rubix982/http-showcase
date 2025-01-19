package main

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/http2"
	"io"
	"log"
	"net/http"
)

func main() {
	// Map to store pushed resources
	pushedResources := make(map[string][]byte)

	// Custom transport for HTTP/2 with a Push handler
	transport := &http2.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("GET", "https://localhost:8443/push", nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// Issue the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Fatalf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	// Read the main response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	fmt.Printf("Main Response: %s\n", string(body))

	// Process pushed resources
	for url, content := range pushedResources {
		fmt.Printf("Pushed Resource URL: %s\n", url)
		fmt.Printf("Content: %s\n", string(content))
	}
}
