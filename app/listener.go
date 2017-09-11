package app

import (
	"log"
	"net"
)

// Listener is used to wrap an underlying listener.
type Listener struct {
	Listener net.Listener
}

func NewListener(listener net.Listener) *Listener {
	return &Listener{
		Listener: listener,
	}
}

// Accept waits for and returns the next connection to the listener.
func (l *Listener) Accept() (net.Conn, error) {
	// Get the underlying connection
	conn, err := l.Listener.Accept()

	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Close closes the underlying listener.
func (l *Listener) Close() error {
	return l.Listener.Close()
}

// Addr returns the underlying listener's network address.
func (l *Listener) Addr() net.Addr {
	return l.Listener.Addr()
}

// Blocks on new connections, and when getting one
// it passes it to the provided channel.
func (l *Listener) awaitForConnections(clientConnChan chan<- net.Conn) {
	for {
		clientConn, err := l.Accept()

		if err != nil {
			log.Println(err)
			continue
		}

		clientConnChan <- clientConn
	}
}
