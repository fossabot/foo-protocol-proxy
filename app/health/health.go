package health

import (
	"net/http"
	"sync"
)

var (
	healthStatus    = http.StatusOK
	readinessStatus = http.StatusOK
	mu              sync.RWMutex
)

func healthzStatus() int {
	mu.RLock()
	defer mu.RUnlock()

	return healthStatus
}

// ReadinessStatus returns current readiness status.
func ReadinessStatus() int {
	mu.RLock()
	defer mu.RUnlock()

	return readinessStatus
}

func setHealthzStatus(status int) {
	mu.Lock()
	healthStatus = status
	mu.Unlock()
}

// SetReadinessStatus sets readiness status to the desired status.
func SetReadinessStatus(status int) {
	mu.Lock()
	readinessStatus = status
	mu.Unlock()
}

// HealthzHandler responds to health check requests.
func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(healthzStatus())
}

// ReadinessHandler responds to readiness check requests.
func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(ReadinessStatus())
}

// ReadinessStatusHandler handles readiness status requests.
func ReadinessStatusHandler(w http.ResponseWriter, r *http.Request) {
	switch ReadinessStatus() {
	case http.StatusOK:
		SetReadinessStatus(http.StatusServiceUnavailable)
	case http.StatusServiceUnavailable:
		SetReadinessStatus(http.StatusOK)
	}
	w.WriteHeader(http.StatusOK)
}

// HealthzStatusHandler handles health status requests.
func HealthzStatusHandler(w http.ResponseWriter, r *http.Request) {
	switch healthzStatus() {
	case http.StatusOK:
		setHealthzStatus(http.StatusServiceUnavailable)
	case http.StatusServiceUnavailable:
		setHealthzStatus(http.StatusOK)
	}
	w.WriteHeader(http.StatusOK)
}
