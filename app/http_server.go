package app

import (
	"foo-protocol-proxy/analysis"
	"foo-protocol-proxy/config"
	"foo-protocol-proxy/handlers"
	"log"
	"net/http"
)

type (
	HttpServer struct {
		config    config.Configuration
		analyzer  *analysis.Analyzer
		errorChan chan error
	}
)

func NewHttpServer(config config.Configuration, analyzer *analysis.Analyzer) *HttpServer {
	return &HttpServer{
		config:    config,
		analyzer:  analyzer,
		errorChan: make(chan error, 10),
	}
}

func (s *HttpServer) Start() {
	s.configureRoutes()

	go func() {
		s.errorChan <- http.ListenAndServe(s.config.HttpAddress, nil)
	}()
	go s.monitorErrors()
}

func (s *HttpServer) configureRoutes() {
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

func (s *HttpServer) getHandler(route string) http.Handler {
	var handler http.Handler

	switch route {
	case "/metrics", "/stats":
		handler = handlers.NewMetricsHandler(s.analyzer)

	case "/health", "/status":
		handler = handlers.NewHealthHandler()
	}

	return handler
}

func (s *HttpServer) monitorErrors() {
	for {
		select {
		case err := <-s.errorChan:
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
