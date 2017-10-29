package app

import (
	"errors"
	"github.com/ahmedkamals/foo-protocol-proxy/analysis"
	"github.com/ahmedkamals/foo-protocol-proxy/config"
	"net/http"
)

type (
	// HTTPServer is an interface for handling HTTP connections.
	HTTPServer struct {
		config    config.Configuration
		routes    map[string]http.Handler
		server    *http.Server
		analyzer  *analysis.Analyzer
		errorChan chan error
	}
)

// NewHTTPServer allocates and returns a new HTTPServer to handle HTTP connections.
func NewHTTPServer(
	config config.Configuration,
	routes map[string]http.Handler,
	errorChan chan error,
) *HTTPServer {
	return &HTTPServer{
		config:    config,
		routes:    routes,
		server:    &http.Server{Addr: config.HTTPAddress},
		errorChan: errorChan,
	}
}

// Start initiates routes configuration, and starts listening.
func (s *HTTPServer) Start() error {
	mux, err := s.configureRoutesHandler(s.routes)

	if err != nil {
		return err
	}

	s.server.Handler = mux

	return s.server.ListenAndServe()
}

func (s *HTTPServer) configureRoutesHandler(routes map[string]http.Handler) (*http.ServeMux, error) {
	if len(routes) == 0 {
		return nil, errors.New("missing routes")
	}
	mux := http.NewServeMux()

	for route, handler := range routes {
		mux.Handle(route, handler)
	}

	return mux, nil
}

func (s *HTTPServer) getRoutes() map[string]http.Handler {
	return s.routes
}

// Close immediately closes all active net.Listeners and any
// connections in state StateNew, StateActive, or StateIdle.
func (s *HTTPServer) Close() error {
	return s.server.Close()
}
