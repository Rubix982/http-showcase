package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"sync"

	"github.com/quic-go/quic-go/http3"
)

func main() {
	// Create a custom HTTP/3 client
	client := &http.Client{
		Transport: &http3.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // For self-signed certificates (not recommended for production)
			},
		},
	}

	baseURL := "https://localhost:8443"

	// URLs to test
	urls := []string{
		baseURL,
		fmt.Sprintf("%s/static/style.css", baseURL),
		fmt.Sprintf("%s/static/script.js", baseURL),
		fmt.Sprintf("%s/push", baseURL),
	}

	for _, url := range urls {
		processRequest(client, url)
	}

	testMultipleStreams(client)
}

func processRequest(client *http.Client, url string) {
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Failed to fetch %s: %v", url, err)
		return // Skip the rest of the function
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Printf("Failed to close response body for %s: %v", url, err)
		}
	}(resp.Body)

	// Read and display the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body for %s: %v", url, err)
		return // Skip the rest of the function
	}

	fmt.Printf("Response from %s:\n", url)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Body:\n%s\n", string(body))
}

func testMultipleStreams(client *http.Client) {
	urls := make([]string, 10)

	for i := 0; i < 10; i++ {
		inRange, err := randomInRange(10, 30)
		if err != nil {
			log.Printf("Failed to generate random number: %v", err)
		}
		urls[i] = fmt.Sprintf("https://localhost:8443/stream%d?delay=%d", i+1, inRange)
	}

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			resp, err := client.Get(url)
			if err != nil {
				log.Printf("Error fetching %s: %v", url, err)
				return
			}

			defer func(Body io.ReadCloser) {
				if err = Body.Close(); err != nil {
					log.Printf("Failed to close response body for %s: %v", url, err)
				}
			}(resp.Body)

			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("Response from %s: %s\n", url, body)
		}(url)
	}
	wg.Wait()
}

// RandomInRange generates a random number between min and max (inclusive)
func randomInRange(min, max int) (int, error) {
	if min > max {
		return 0, fmt.Errorf("invalid range: min (%d) cannot be greater than max (%d)", min, max)
	}

	// Calculate the range size
	rangeSize := max - min + 1

	// Generate a random number in the range [0, rangeSize)
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(rangeSize)))
	if err != nil {
		return 0, fmt.Errorf("failed to generate random number: %v", err)
	}

	// Adjust to the desired range
	return int(nBig.Int64()) + min, nil
}
