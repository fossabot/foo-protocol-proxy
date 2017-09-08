package app

import (
	"flag"
	"foo-protocol-proxy/config"
)

type (
	Dispatcher struct {
	}
)

func (d *Dispatcher) Run() {
	config := d.parseConfig()
	NewProxy(config).Start()
}

func (d *Dispatcher) parseConfig() config.Configuration {
	listen := flag.String("listen", ":8002", "Listening port.")
	forward := flag.String("forward", ":8001", "Forwarding port.")
	flag.Parse()

	return config.Configuration{
		Listening:  *listen,
		Forwarding: *forward,
	}
}
