package main

import (
	"foo-protocol-proxy/app"
)

func main() {
	dispatcher := new(app.Dispatcher)
	dispatcher.Run()
}
