package xrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Handler is a function that handles an XRPC request
type Handler func(ctx context.Context, params map[string]interface{}) (interface{}, error)

// HandlerMap maps procedure names to handlers
type HandlerMap map[string]Handler

// Server is an XRPC server
type Server struct {
	handlers HandlerMap
}

// NewServer creates a new XRPC server
func NewServer() *Server {
	return &Server{
		handlers: make(HandlerMap),
	}
}

// Register registers a handler for a procedure
func (s *Server) Register(procedure string, handler Handler) {
	s.handlers[procedure] = handler
}

// ServeHTTP implements the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract procedure name from path
	path := strings.TrimPrefix(r.URL.Path, "/xrpc/")
	procedure := strings.TrimSuffix(path, "/")

	// Find handler for procedure
	handler, ok := s.handlers[procedure]
	if !ok {
		http.Error(w, fmt.Sprintf("unknown procedure: %s", procedure), http.StatusNotFound)
		return
	}

	// Parse parameters
	params := make(map[string]interface{})
	
	// Parse query parameters
	for key, values := range r.URL.Query() {
		if len(values) == 1 {
			params[key] = values[0]
		} else {
			params[key] = values
		}
	}

	// Parse body if present
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		var bodyParams map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&bodyParams); err == nil {
			for key, value := range bodyParams {
				params[key] = value
			}
		}
	}

	// Execute handler
	result, err := handler(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Register registers a handler with an http.ServeMux
func Register(mux *http.ServeMux, procedure string, handler Handler) {
	server := NewServer()
	server.Register(procedure, handler)
	mux.Handle("/xrpc/"+procedure+"/", server)
}
