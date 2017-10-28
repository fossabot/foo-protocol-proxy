package app

import (
	"fmt"
	"github.com/ahmedkamals/foo-protocol-proxy/analysis"
	"github.com/ahmedkamals/foo-protocol-proxy/config"
	"github.com/ahmedkamals/foo-protocol-proxy/persistence"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type (
	// Dispatcher acts as en entry point for the application.
	Dispatcher struct {
		config   config.Configuration
		analyzer *analysis.Analyzer
		saver    *persistence.Saver
		proxy    *Proxy
	}
)

// NewDispatcher allocates and returns a new Dispatcher.
func NewDispatcher(config config.Configuration, analyzer *analysis.Analyzer, saver *persistence.Saver) *Dispatcher {
	return &Dispatcher{
		config:   config,
		analyzer: analyzer,
		saver:    saver,
	}
}

// Run starts the dispatcher.
func (d *Dispatcher) Run() {
	d.proxy = NewProxy(d.config, d.analyzer, d.saver)
	err := d.proxy.Start()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	NewHTTPServer(d.config, d.analyzer).Start()

	if d.blockIndefinitely(make(chan os.Signal, 1), true) {
		d.Close()
	}
}

// Close closes the dispatcher and its dependencies.
func (d *Dispatcher) Close() {
	if d.proxy != nil {
		d.proxy.Close()
	}
	os.Exit(0)
}

// blockIndefinitely blocks for interrupt signal from the OS.
func (d *Dispatcher) blockIndefinitely(signalChan chan os.Signal, breakOnSignal bool) bool {
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	for {
		select {
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exiting...", s))

			if breakOnSignal {
				return true
			}
		}
	}
}
