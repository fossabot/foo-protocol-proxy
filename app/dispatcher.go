package app

import (
	"github.com/ahmedkamals/foo-protocol-proxy/analysis"
	"github.com/ahmedkamals/foo-protocol-proxy/config"
	"github.com/ahmedkamals/foo-protocol-proxy/handlers"
	"github.com/ahmedkamals/foo-protocol-proxy/persistence"
	"github.com/kpango/glg"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type (
	// Dispatcher acts as en entry point for the application.
	Dispatcher struct {
		config     config.Configuration
		proxy      *Proxy
		httpServer *HTTPServer
		analyzer   *analysis.Analyzer
		saver      persistence.Saver
		logger     *glg.Glg
		errorChan  chan error
	}
)

// NewDispatcher allocates and returns a new Dispatcher.
func NewDispatcher(config config.Configuration, analyzer *analysis.Analyzer, saver persistence.Saver) *Dispatcher {
	return &Dispatcher{
		config:    config,
		analyzer:  analyzer,
		saver:     saver,
		logger:    glg.New(),
		errorChan: make(chan error, 10),
	}
}

// Start starts the dispatcher.
func (d *Dispatcher) Start() {
	d.proxy = NewProxy(d.config, d.analyzer, d.saver, d.logger, d.errorChan)
	err := d.proxy.Start()

	if err != nil {
		d.logger.Error(err)
		os.Exit(1)
	}

	d.httpServer = NewHTTPServer(d.config, d.getRoutes(), d.errorChan)

	go func() {
		d.errorChan <- d.httpServer.Start()
	}()

	if d.blockIndefinitely(make(chan os.Signal, 1), true) {
		d.Close()
	}
}

// Close closes the dispatcher and its dependencies.
func (d *Dispatcher) Close() {
	if d.proxy != nil {
		d.proxy.Close()
		d.httpServer.Close()
	}
}

// blockIndefinitely blocks for interrupt signal from the OS.
func (d *Dispatcher) blockIndefinitely(signalChan chan os.Signal, breakOnSignal bool) bool {
	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	for s := range signalChan {

		d.logger.Infof("Captured %v. Exiting...", s)

		if breakOnSignal {
			return true
		}
	}

	return false
}

func (d *Dispatcher) getRoutes() map[string]http.Handler {
	eventLogger := log.New(os.Stdout, "", 0)

	return map[string]http.Handler{
		"/metrics":   handlers.NewMetricsHandler(d.analyzer, eventLogger),
		"/stats":     handlers.NewMetricsHandler(d.analyzer, eventLogger),
		"/health":    handlers.NewHealthHandler(eventLogger),
		"/heartbeat": handlers.NewHealthHandler(eventLogger),
		"/status":    handlers.NewHealthHandler(eventLogger),
	}
}
