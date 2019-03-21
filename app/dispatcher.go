package app

import (
	"github.com/ahmedkamals/foo-protocol-proxy/app/analysis"
	"github.com/ahmedkamals/foo-protocol-proxy/app/config"
	"github.com/ahmedkamals/foo-protocol-proxy/app/handlers"
	"github.com/ahmedkamals/foo-protocol-proxy/app/health"
	"github.com/ahmedkamals/foo-protocol-proxy/app/persistence"
	"github.com/ahmedkamals/foo-protocol-proxy/app/pkg/version"
	"github.com/braintree/manners"
	"github.com/kpango/glg"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type (
	// Dispatcher acts as en entry point for the application.
	Dispatcher struct {
		config       config.Configuration
		proxy        *Proxy
		healthServer *HTTPServer
		httpServer   *HTTPServer
		analyzer     *analysis.Analyzer
		saver        persistence.Saver
		logger       *glg.Glg
		errChan      chan error
	}
)

// NewDispatcher allocates and returns a new Dispatcher.
func NewDispatcher(config config.Configuration, analyzer *analysis.Analyzer, saver persistence.Saver) *Dispatcher {
	return &Dispatcher{
		config:   config,
		analyzer: analyzer,
		saver:    saver,
		logger:   glg.New(),
		errChan:  make(chan error, 10),
	}
}

// Start starts the dispatcher.
func (d *Dispatcher) Start() {
	d.errChan <- d.logger.Info("Starting server...")
	d.proxy = NewProxy(d.config, d.analyzer, d.saver, d.logger, d.errChan)
	err := d.proxy.Start()

	if err != nil {
		d.errChan <- d.logger.Error(err)
		os.Exit(1)
	}

	d.startHealthServer()
	d.httpServer = NewHTTPServer(d.config, d.getRoutes(), d.errChan)
	go d.monitorErrors()

	go func() {
		d.errChan <- d.logger.Infof("HTTP service listening on %s", d.config.HTTPAddress)
		d.errChan <- d.httpServer.Start()
	}()

	if d.blockIndefinitely(make(chan os.Signal, 1), true) {
		d.Close()
	}
}

// Close closes the dispatcher and its dependencies.
func (d *Dispatcher) Close() {
	if d.proxy != nil {
		d.proxy.Close()
	}
	health.SetReadinessStatus(http.StatusServiceUnavailable)
	d.errChan <- d.httpServer.Close()
	close(d.errChan)
	os.Exit(0)
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
		d.errChan <- d.logger.Infof("Captured %v. Exiting...", s)

		if breakOnSignal {
			return true
		}
	}

	return false
}

func (d *Dispatcher) getRoutes() map[string]http.Handler {
	eventLogger := log.New(os.Stdout, "", 0)

	return map[string]http.Handler{
		"/metrics": handlers.MetricsHandler(d.analyzer, eventLogger),
		"/stats":   handlers.MetricsHandler(d.analyzer, eventLogger),
		"/version": handlers.VersionHandler(version.VERSION, eventLogger),
	}
}

func (d *Dispatcher) startHealthServer() {
	hmux := http.NewServeMux()
	hmux.HandleFunc("/health", health.HealthzHandler)
	hmux.HandleFunc("/heartbeat", health.HealthzHandler)
	hmux.HandleFunc("/status", health.HealthzHandler)
	hmux.HandleFunc("/healthz/status", health.HealthzStatusHandler)
	hmux.HandleFunc("/readiness", health.ReadinessHandler)
	hmux.HandleFunc("/readiness/status", health.ReadinessStatusHandler)

	healthServer := manners.NewServer()
	healthServer.Addr = d.config.HealthAddress
	healthServer.Handler = handlers.LoggingHandler(hmux)

	go func() {
		d.errChan <- d.logger.Infof("Health service listening on %s", d.config.HealthAddress)
		d.errChan <- healthServer.ListenAndServe()
	}()
}

func (d *Dispatcher) monitorErrors() {
	for err := range d.errChan {
		if err != nil && err != io.EOF {
			d.logger.Error(err)
		}
	}
}
