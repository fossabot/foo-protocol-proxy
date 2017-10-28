package app

import (
	"github.com/ahmedkamals/foo-protocol-proxy/analysis"
	"github.com/ahmedkamals/foo-protocol-proxy/config"
	"github.com/ahmedkamals/foo-protocol-proxy/handlers"
	"log"
	"net/http"
)

type (
	// HTTPServer is an interface for handling HTTP connections.
	HTTPServer struct {
		config    config.Configuration
		analyzer  *analysis.Analyzer
		errorChan chan error
	}
)

// NewHTTPServer allocates and returns a new HTTPServer to handle HTTP connections.
func NewHTTPServer(config config.Configuration, analyzer *analysis.Analyzer) *HTTPServer {
	return &HTTPServer{
		config:    config,
		analyzer:  analyzer,
		errorChan: make(chan error, 10),
	}
}

// Start initiates routes configuration, and starts listening.
func (s *HTTPServer) Start() {
	routes := s.getRoutes()
	s.configureRoutes(routes)

	go func() {
		s.errorChan <- http.ListenAndServe(s.config.HTTPAddress, nil)
	}()
	go s.monitorErrors()
}

func (s *HTTPServer) configureRoutes(routes map[string]http.Handler) {
	for route, handler := range routes {
		http.Handle(route, handler)
	}
}

func (s *HTTPServer) getRoutes() map[string]http.Handler {
	return map[string]http.Handler{
		"/metrics": handlers.NewMetricsHandler(s.analyzer),
		"/stats":   handlers.NewMetricsHandler(s.analyzer),
		"/health":  handlers.NewHealthHandler(),
		"/status":  handlers.NewHealthHandler(),
	}
}

func (s *HTTPServer) monitorErrors() {
	for {
		select {
		case err := <-s.errorChan:
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
