package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Service represents a microservice
type Service struct {
	Name string
	URL  string
}

// Services map
var services = map[string]Service{
	"auth": {Name: "Authentication Service", URL: "http://localhost:8081"},
	"pds":  {Name: "Personal Data Server", URL: "http://localhost:8082"},
	"bgs":  {Name: "Big Graph Service", URL: "http://localhost:8083"},
}

// ProxyHandler forwards requests to the appropriate service
func ProxyHandler(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the URL
		remote, err := url.Parse(service.URL)
		if err != nil {
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			log.Printf("Error parsing service URL: %v", err)
			return
		}

		// Create the reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(remote)
		
		// Update the headers to allow for SSL redirection
		r.URL.Host = remote.Host
		r.URL.Scheme = remote.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		
		// Note that ServeHttp is non-blocking
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	// Create HTTP server
	mux := http.NewServeMux()

	// Auth service routes
	mux.HandleFunc("/auth/", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/auth")
		ProxyHandler(services["auth"])(w, r)
	})

	// PDS routes
	mux.HandleFunc("/pds/", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/pds")
		ProxyHandler(services["pds"])(w, r)
	})

	// BGS routes
	mux.HandleFunc("/bgs/", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/bgs")
		ProxyHandler(services["bgs"])(w, r)
	})

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Create server
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting API Gateway on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Graceful shutdown
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")
}
