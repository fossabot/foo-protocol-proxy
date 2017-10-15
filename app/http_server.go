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
	s.configureRoutes()

	go func() {
		s.errorChan <- http.ListenAndServe(s.config.HTTPAddress, nil)
	}()
	go s.monitorErrors()
}

func (s *HTTPServer) configureRoutes() {
	routes := []string{
		"/metrics",
		"/stats",
		"/health",
		"/status",
	}

	for _, route := range routes {
		http.Handle(route, s.getHandler(route))
	}
}

func (s *HTTPServer) getHandler(route string) http.Handler {
	var handler http.Handler

	switch route {
	case "/metrics", "/stats":
		handler = handlers.NewMetricsHandler(s.analyzer)

	case "/health", "/status":
		handler = handlers.NewHealthHandler()
	}

	return handler
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
