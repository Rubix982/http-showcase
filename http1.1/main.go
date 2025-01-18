package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
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
	if r.Header.Get("Connection") == "keep-alive" {
		w.Header().Set("Connection", "keep-alive")
	}
	_, _ = fmt.Fprintln(w, "Welcome to the HTTP/1.1 Server!")
}

func chunkedHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)

	chunks := []string{"Chunk 1: Hello", "Chunk 2: HTTP/1.1", "Chunk 3: Chunked response"}
	for _, chunk := range chunks {
		time.Sleep(1 * time.Second)
		_, _ = io.WriteString(w, fmt.Sprintf("%x\r\n%s\r\n", len(chunk), chunk))
	}
	_, _ = io.WriteString(w, "0\r\n\r\n")
}

func pipeliningHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	_, _ = fmt.Fprintf(w, "Pipelined response for %s\n", r.URL.Path)
}

func connectionStateHandler(_ net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew:
		atomic.AddInt32(&connectionCounter, 1)
		log.Printf("New connection established. Active connections: %d", atomic.LoadInt32(&connectionCounter))
	case http.StateIdle:
		log.Printf("Connection is idle. Active connections: %d", atomic.LoadInt32(&connectionCounter))
	case http.StateClosed:
		atomic.AddInt32(&connectionCounter, -1)
		log.Printf("Connection closed. Active connections: %d", atomic.LoadInt32(&connectionCounter))
	case http.StateActive:
		log.Printf("Connection is active. Active connections: %d", atomic.LoadInt32(&connectionCounter))
	case http.StateHijacked:
		log.Printf("Connection hijacked. Active connections: %d", atomic.LoadInt32(&connectionCounter))
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/chunked", chunkedHandler)
	mux.HandleFunc("/pipelining", pipeliningHandler)

	loggedMux := loggingMiddleware(mux)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           loggedMux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		MaxHeaderBytes:    1 << 20,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		ConnState:         connectionStateHandler,
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

	fmt.Println("HTTP/1.1 Server is running on port 8080")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Error starting server: %v", err)
	}
}
