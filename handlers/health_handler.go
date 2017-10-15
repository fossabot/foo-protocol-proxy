package handlers

import (
	"encoding/json"
	"net/http"
)

type (
	// HealthHandler acts as an interface for server health over HTTP.
	HealthHandler struct {
	}
)

// NewHealthHandler allocates and returns a new HealthHandler to report health.
func NewHealthHandler() http.Handler {
	return new(HealthHandler)
}

func (m *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType != "" {
		contentType = "application/json"
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(map[string]string{"status": "OK!"})
}
