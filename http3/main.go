package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	rand2 "math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/quic-go/quic-go/http3"
)

// streamHandler creates an HTTP handler for a specific stream ID
// and processes a delay provided as a query parameter (?delay=ms).
func streamHandler(id int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the delay query parameter
		delayStr := r.URL.Query().Get("delay")
		delay, err := strconv.Atoi(delayStr)
		if err != nil || delay < 0 {
			http.Error(w, "Invalid or missing delay parameter", http.StatusBadRequest)
			return
		}

		// Simulate processing delay
		time.Sleep(time.Duration(delay) * time.Second)

		// Respond with a message indicating the stream ID and delay
		response := fmt.Sprintf("Stream %d responded after %d ms", id, delay)
		log.Println(response)
	}
}

func artificialPacketLoss(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rand2.Float64() < 0.2 { // 20% packet loss
			log.Println("Simulating packet loss")
			return // Drop the request
		}
		next.ServeHTTP(w, r)
	})
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Path
	if file == "/" {
		file = "/index.html"
	}
	http.ServeFile(w, r, "static"+file)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func pushHandler(w http.ResponseWriter, _ *http.Request) {
	// Example of server push for resources
	if pusher, ok := w.(http.Pusher); ok {
		filesToPush := []string{"/style.css", "/script.js"}
		for _, file := range filesToPush {
			if err := pusher.Push(file, nil); err != nil {
				log.Printf("Failed to push %s: %v", file, err)
			}
		}
	}

	// Serve the main response
	log.Print("Hello from HTTP/3 with Server Push!")
}

func main() {
	// Serve static files
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	for i := 1; i <= 10; i++ {
		mux.HandleFunc(fmt.Sprintf("/stream%d", i), streamHandler(i))
	}
	mux.HandleFunc("/", staticHandler)
	mux.HandleFunc("/push", pushHandler)

	loggedMux := loggingMiddleware(mux)
	artificialMux := artificialPacketLoss(loggedMux)

	// Generate or load TLS certificates
	certFile, keyFile := "server.crt", "server.key"
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		log.Println("Generating self-signed certificates...")
		if err := generateSelfSignedCertificate(certFile, keyFile); err != nil {
			log.Fatalf("Failed to generate certificates: %v", err)
		}
	}

	// Setup HTTP/3 server
	server := &http3.Server{
		Addr:      ":8443",
		Handler:   artificialMux,
		TLSConfig: generateTLSConfig(certFile, keyFile),
	}

	log.Println("HTTP/3 Server is running on https://localhost:8443")
	if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// generateTLSConfig loads or generates the required TLS configuration for QUIC
func generateTLSConfig(certFile, keyFile string) *tls.Config {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load TLS certificates: %v", err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h3"}, // Enable HTTP/3
	}
}

// generateSelfSignedCert generates a self-signed certificate (for testing purposes)
func generateSelfSignedCertificate(certFile, keyFile string) error {
	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"HTTP/3 Server"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	certOut, err := os.Create(certFile)
	if err != nil {
		return err
	}

	defer func(certOut *os.File) {
		if err = certOut.Close(); err != nil {
			log.Fatalf("Failed to close cert file: %v", err)
		}
	}(certOut)

	if err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return err
	}

	keyOut, err := os.Create(keyFile)
	if err != nil {
		return err
	}

	defer func(keyOut *os.File) {
		if err = keyOut.Close(); err != nil {
			log.Fatalf("Failed to close key file: %v", err)
		}
	}(keyOut)

	keyBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return err
	}
	if err = pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes}); err != nil {
		return err
	}

	return nil
}
