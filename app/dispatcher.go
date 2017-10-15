package app

import (
	"flag"
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
		proxy *Proxy
	}
)

// Run starts the dispatcher.
func (d *Dispatcher) Run() {
	config := d.parseConfig()
	analyzer := analysis.NewAnalyzer()
	saver := persistence.NewSaver(config.RecoveryPath)

	d.proxy = NewProxy(config, analyzer, saver)
	err := d.proxy.Start()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	NewHTTPServer(config, analyzer).Start()

	d.blockIndefinitely()
}

func (d *Dispatcher) parseConfig() config.Configuration {
	var (
		listen       = flag.String("listen", ":8002", "Listening port.")
		forward      = flag.String("forward", ":8001", "Forwarding port.")
		httpAddr     = flag.String("http", "0.0.0.0:8088", "Health service address.")
		recoveryPath = flag.String("recovery-path", "data/recovery.json", "Recovery path.")
	)
	flag.Parse()

	return config.Configuration{
		Listening:    *listen,
		Forwarding:   *forward,
		HTTPAddress:  *httpAddr,
		RecoveryPath: *recoveryPath,
	}
}

func (d *Dispatcher) blockIndefinitely() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	for {
		select {
		case s := <-signalChan:
			d.proxy.close()
			log.Println(fmt.Sprintf("Captured %v. Exiting...", s))
			os.Exit(0)
		}
	}
}
