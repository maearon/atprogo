package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/yourusername/atprogo/pkg/bgs"
	"github.com/yourusername/atprogo/pkg/db"
)

// FollowRequest represents a follow request
type FollowRequest struct {
	Follower  string `json:"follower"`
	Following string `json:"following"`
}

// BGSHandler handles BGS requests
type BGSHandler struct {
	followRepo *bgs.FollowRepository
}

// NewBGSHandler creates a new BGS handler
func NewBGSHandler(followRepo *bgs.FollowRepository) *BGSHandler {
	return &BGSHandler{
		followRepo: followRepo,
	}
}

// FollowHandler handles follow requests
func (h *BGSHandler) FollowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Follower == "" || req.Following == "" {
		http.Error(w, "Follower and following are required", http.StatusBadRequest)
		return
	}

	// Create follow
	follow := &bgs.Follow{
		Follower:  req.Follower,
		Following: req.Following,
	}
	if err := h.followRepo.CreateFollow(r.Context(), follow); err != nil {
		log.Printf("Failed to create follow: %v", err)
		http.Error(w, "Failed to create follow", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(follow)
}

// UnfollowHandler handles unfollow requests
func (h *BGSHandler) UnfollowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Follower == "" || req.Following == "" {
		http.Error(w, "Follower and following are required", http.StatusBadRequest)
		return
	}

	// Delete follow
	if err := h.followRepo.DeleteFollow(r.Context(), req.Follower, req.Following); err != nil {
		log.Printf("Failed to delete follow: %v", err)
		http.Error(w, "Failed to delete follow", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// GetFollowersHandler handles getting followers
func (h *BGSHandler) GetFollowersHandler(w http.ResponseWriter, r *http.Request) {
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

	// Get limit and offset from query parameters
	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get followers
	followers, err := h.followRepo.GetFollowers(r.Context(), did, limit, offset)
	if err != nil {
		log.Printf("Failed to get followers: %v", err)
		http.Error(w, "Failed to get followers", http.StatusInternalServerError)
		return
	}

	// Get followers count
	count, err := h.followRepo.GetFollowersCount(r.Context(), did)
	if err != nil {
		log.Printf("Failed to get followers count: %v", err)
		http.Error(w, "Failed to get followers count", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"followers": followers,
		"count":     count,
	})
}

// GetFollowingHandler handles getting following
func (h *BGSHandler) GetFollowingHandler(w http.ResponseWriter, r *http.Request) {
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

	// Get limit and offset from query parameters
	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get following
	following, err := h.followRepo.GetFollowing(r.Context(), did, limit, offset)
	if err != nil {
		log.Printf("Failed to get following: %v", err)
		http.Error(w, "Failed to get following", http.StatusInternalServerError)
		return
	}

	// Get following count
	count, err := h.followRepo.GetFollowingCount(r.Context(), did)
	if err != nil {
		log.Printf("Failed to get following count: %v", err)
		http.Error(w, "Failed to get following count", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"following": following,
		"count":     count,
	})
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
	followRepo := bgs.NewFollowRepository(dbPool)

	// Create handlers
	bgsHandler := NewBGSHandler(followRepo)

	// Create HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/follow", bgsHandler.FollowHandler)
	mux.HandleFunc("/unfollow", bgsHandler.UnfollowHandler)
	mux.HandleFunc("/followers", bgsHandler.GetFollowersHandler)
	mux.HandleFunc("/following", bgsHandler.GetFollowingHandler)

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	// Create server
	server := &http.Server{
		Addr:    ":8083",
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting BGS service on port 8083")
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
