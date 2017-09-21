package app

import (
	"flag"
	"fmt"
	"github.com/ahmedkamals/foo-protocol-proxy/analysis"
	"github.com/ahmedkamals/foo-protocol-proxy/config"
	"github.com/ahmedkamals/foo-protocol-proxy/persistance"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type (
	Dispatcher struct {
		proxy *Proxy
	}
)

func (d *Dispatcher) Run() {
	config := d.parseConfig()
	analyzer := analysis.NewAnalyzer()
	saver := persistance.NewSaver(config.RecoveryPath)

	d.proxy = NewProxy(config, analyzer, saver)
	err := d.proxy.Start()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	NewHttpServer(config, analyzer).Start()

	d.BlockIndefinitely()
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
		HttpAddress:  *httpAddr,
		RecoveryPath: *recoveryPath,
	}
}

func (d *Dispatcher) BlockIndefinitely() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	for {
		select {
		case s := <-signalChan:
			d.proxy.Close()
			log.Println(fmt.Sprintf("Captured %v. Exiting...", s))
			os.Exit(0)
		}
	}
}
