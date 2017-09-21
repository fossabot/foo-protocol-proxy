package main

import (
	"github.com/ahmedkamals/foo-protocol-proxy/app"
)

func main() {
	dispatcher := new(app.Dispatcher)
	dispatcher.Run()
}
