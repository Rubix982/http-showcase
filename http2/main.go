package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/net/http2"
)

var connectionCounter int32

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Host == "" {
		http.Error(w, "400 Bad Request - Missing Host Header", http.StatusBadRequest)
		return
	}

	if w.Header().Get("Connection") == "keep-alive" {
		w.Header().Set("Connection", "keep-alive")
	}

	log.Print("Welcome to the HTTP/1.1 Server!")
}

// Multiplexed streams handler
func multiplexedHandler(_ http.ResponseWriter, _ *http.Request) {
	// Simulate multiple streams by sending multiple message over the same connection
	for i := 1; i <= 5; i++ {
		log.Printf("Message %d from streams\n", i)
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

	log.Print("Headers added. Check with an HTTP/2 client.")
}

// Stream priority handler
func prioritizationHandler(_ http.ResponseWriter, r *http.Request) {
	priority := r.URL.Query().Get("priority")
	switch priority {
	case "high":
		time.Sleep(1 * time.Second)
		log.Print("High priority stream")
	case "low":
		time.Sleep(3 * time.Second)
		log.Print("Low priority stream")
	default:
		time.Sleep(2 * time.Second)
		log.Print("Default priority stream")
	}
	log.Print("Stream priority demonstration. Check with an HTTP/2 client.")
}

// Server push handler
// !!!!!!!!!
// This logic is currently not working. Accessing it via the browser provides the error that the feature is not supported.
// Using curl with `curl -k --http2 -v https://localhost:8443/push` shows that the feature is not supported.
// Either it is because this feature is not supported or adopted well. I cannot find any other reason for this.
func serverPushHandler(w http.ResponseWriter, r *http.Request) {
	if pusher, ok := w.(http.Pusher); ok {
		pushOptions := http.PushOptions{
			Method: r.Method,
			Header: r.Header,
		}
		for _, file := range []string{"static/style.css", "static/script.js"} {
			_, err := os.Stat(file)
			if err != nil {
				log.Fatalf("%s not found: %v", file, err)
			}
			if !strings.HasPrefix(file, "/") {
				file = fmt.Sprintf("/%v", file)
			}
			if err = pusher.Push(file, &pushOptions); err != nil {
				log.Printf("Failed to push %s. Error: %v", file, err)
			}
		}
	}
	log.Print("Stream priority demonstration. Check for pushed resources. Check with an HTTP/2 client.")
}

func connectionStateHandler(networkConnection net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew:
		atomic.AddInt32(&connectionCounter, 1)
		log.Printf("New connection established.")
	case http.StateIdle:
		log.Printf("Connection is idle.")
	case http.StateClosed:
		atomic.AddInt32(&connectionCounter, -1)
		log.Printf("Connection closed.")
	case http.StateActive:
		log.Printf("Connection is active.")
	case http.StateHijacked:
		log.Printf("Connection hijacked.")
	}
	log.Printf("Active connections: %v. Remote address: %v.", atomic.LoadInt32(&connectionCounter), networkConnection.RemoteAddr())
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/multiplex", multiplexedHandler)
	mux.HandleFunc("/headers", headerCompressionHandler)
	mux.HandleFunc("/priority", prioritizationHandler)
	mux.HandleFunc("/push", serverPushHandler)

	loggedMux := loggingMiddleware(mux)

	server := &http.Server{
		Addr:      ":8443",
		Handler:   loggedMux,
		ConnState: connectionStateHandler,
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
