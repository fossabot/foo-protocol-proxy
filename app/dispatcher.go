package app

import (
	"flag"
	"fmt"
	"foo-protocol-proxy/analysis"
	"foo-protocol-proxy/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type (
	Dispatcher struct {
	}
)

func (d *Dispatcher) Run() {
	config := d.parseConfig()
	analyzer := analysis.NewAnalyzer()
	err := NewProxy(config, analyzer).Start()
	NewHttpServer(config, analyzer).Start()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	d.BlockIndefinitely()
}

func (d *Dispatcher) parseConfig() config.Configuration {
	var (
		listen   = flag.String("listen", ":8002", "Listening port.")
		forward  = flag.String("forward", ":8001", "Forwarding port.")
		httpAddr = flag.String("http", "0.0.0.0:8088", "Health service address.")
	)
	flag.Parse()

	return config.Configuration{
		Listening:   *listen,
		Forwarding:  *forward,
		HttpAddress: *httpAddr,
	}
}

func (*Dispatcher) BlockIndefinitely() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	for {
		select {
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exiting...", s))
			os.Exit(0)
		}
	}
}
