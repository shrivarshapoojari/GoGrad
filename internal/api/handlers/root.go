package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// RootHandler handles requests to the root endpoint "/"
func RootHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method for root endpoint
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create a welcome response with API information
	response := struct {
		Status    string   `json:"status"`
		Message   string   `json:"message"`
		Version   string   `json:"version"`
		Timestamp string   `json:"timestamp"`
		Endpoints []string `json:"available_endpoints"`
	}{
		Status:    "success",
		Message:   "Welcome to GoGrad REST API",
		Version:   "v1.0.0",
		Timestamp: time.Now().Format(time.RFC3339),
		Endpoints: []string{
			"GET /",
			"GET /teachers",
			"POST /teachers",
			"GET /teachers/{id}",
			"GET /students",
			"POST /students",
			"GET /students/{id}",
		},
	}

	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode and send the response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// HealthHandler provides a simple health check endpoint
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := struct {
		Status    string `json:"status"`
		Message   string `json:"message"`
		Timestamp string `json:"timestamp"`
	}{
		Status:    "healthy",
		Message:   "Server is running successfully",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
