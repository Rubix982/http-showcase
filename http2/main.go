package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/net/http2"
)

// Multiplexed streams handler
func multiplexedHandler(w http.ResponseWriter, _ *http.Request) {
	// Simulate multiple streams by sending multiple message over the same connection
	for i := 1; i <= 5; i++ {
		_, _ = fmt.Fprintf(w, "Message %d from streams\n", i)
		time.Sleep(500 * time.Millisecond)
	}
}

// Header compression (HPACK) demonstration handler
func headerCompressionHandler(w http.ResponseWriter, _ *http.Request) {
	// Add multiple headers to the response to demonstrate header compression
	w.Header().Set("X-Header-1", "Value 1")
	w.Header().Set("X-Header-2", "Value 2")
	w.Header().Set("X-Header-3", "Value 3")
	w.Header().Set("X-Header-4", "Value 4")
	w.Header().Set("X-Header-5", "Value 5")
	w.Header().Set("X-Header-6", "Value 6")
	w.Header().Set("X-Header-7", "Value 7")
	w.Header().Set("X-Header-8", "Value 8")
	w.Header().Set("X-Header-9", "Value 9")
	w.Header().Set("X-Header-10", "Value 10")

	_, _ = fmt.Fprintln(w, "Headers added. Check with an HTTP/2 client.")
}

// Stream priority handler
func prioritizationHandler(w http.ResponseWriter, r *http.Request) {
	priority := r.URL.Query().Get("priority")
	switch priority {
	case "high":
		time.Sleep(1 * time.Second)
		_, _ = fmt.Fprintln(w, "High priority stream")
	case "low":
		time.Sleep(3 * time.Second)
		_, _ = fmt.Fprintln(w, "Low priority stream")
	default:
		time.Sleep(2 * time.Second)
		_, _ = fmt.Fprintln(w, "Default priority stream")
	}
	_, _ = fmt.Fprintln(w, "Stream priority demonstration. Check with an HTTP/2 client.")
}

// Server push handler
func serverPushHandler(w http.ResponseWriter, _ *http.Request) {
	if pusher, ok := w.(http.Pusher); ok {
		// Push related resources
		if err := pusher.Push("https://localhost:8443/style.css", nil); err != nil {
			fmt.Println("Failed to push style.css. Error: ", err)
		}
		if err := pusher.Push("https://localhost:8443/script.js", nil); err != nil {
			fmt.Println("Failed to push script.js. Error: ", err)
		}
	}
	_, _ = fmt.Fprintln(w, "Stream priority demonstration. Check for pushed resources. Check with an HTTP/2 client.")
}

func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/style.css", fs)
	mux.Handle("/script.js", fs)

	mux.HandleFunc("/multiplex", multiplexedHandler)
	mux.HandleFunc("/headers", headerCompressionHandler)
	mux.HandleFunc("/priority", prioritizationHandler)
	mux.HandleFunc("/push", serverPushHandler)

	server := &http.Server{
		Addr:    ":8443",
		Handler: mux,
	}

	h2Config := &http2.Server{}
	if err := http2.ConfigureServer(server, h2Config); err != nil {
		log.Fatalf("Failed to configure HTTP/2 server: %v", err)
	}

	// Graceful shutdown setup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go func() {
		<-stop
		log.Println("Shutting down server gracefully...")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("Error during shutdown: %v", err)
		}
	}()

	log.Println("HTTP/2 Server is running on https://localhost:8443")
	if err := server.ListenAndServeTLS("server.crt", "server.key"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server error: %v", err)
	}
}
