package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/atprogo/pkg/db"
	"github.com/yourusername/atprogo/pkg/pds"
)

// CreatePostRequest represents a create post request
type CreatePostRequest struct {
	DID     string `json:"did"`
	Content string `json:"content"`
}

// GetPostsRequest represents a get posts request
type GetPostsRequest struct {
	DID string `json:"did"`
}

// PDSHandler handles PDS requests
type PDSHandler struct {
	repoRepo *pds.RepositoryRepository
}

// NewPDSHandler creates a new PDS handler
func NewPDSHandler(repoRepo *pds.RepositoryRepository) *PDSHandler {
	return &PDSHandler{
		repoRepo: repoRepo,
	}
}

// CreatePostHandler handles post creation
func (h *PDSHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.DID == "" || req.Content == "" {
		http.Error(w, "DID and content are required", http.StatusBadRequest)
		return
	}

	// Create post
	doc, err := h.repoRepo.CreatePost(r.Context(), req.DID, req.Content)
	if err != nil {
		log.Printf("Failed to create post: %v", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(doc)
}

// GetPostsHandler handles getting posts
func (h *PDSHandler) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get DID from query parameter
	did := r.URL.Query().Get("did")
	if did == "" {
		http.Error(w, "DID is required", http.StatusBadRequest)
		return
	}

	// Get posts
	docs, err := h.repoRepo.GetDocumentsByType(r.Context(), did, "app.bsky.feed.post")
	if err != nil {
		log.Printf("Failed to get posts: %v", err)
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}

func main() {
	// Create context
	ctx := context.Background()

	// Connect to database
	dbPool, err := db.GetDBPool(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Create repositories
	repoRepo := pds.NewRepositoryRepository(dbPool)

	// Create handlers
	pdsHandler := NewPDSHandler(repoRepo)

	// Create HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/posts/create", pdsHandler.CreatePostHandler)
	mux.HandleFunc("/posts/get", pdsHandler.GetPostsHandler)

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	// Create server
	server := &http.Server{
		Addr:    ":8082",
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting PDS service on port 8082")
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
