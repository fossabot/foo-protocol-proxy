package communication

import (
	"bufio"
	"github.com/ahmedkamals/foo-protocol-proxy/analysis"
	"net"
)

// BridgeConnection is used to wrap underlying connections
// with server and client.
type BridgeConnection struct {
	// Source connection
	srcConn net.Conn
	// Destination connection
	dstConn net.Conn
	// Analysis channel
	analysisChan analysis.Aggregation
	// Errors chanel
	errChan chan error
}

// NewBridgeConnection returns a new BridgeConnection
// that forwards connections from source to destination and vice versa.
func NewBridgeConnection(
	srcConn,
	dstConn net.Conn,
	analysisChan analysis.Aggregation,
	errChan chan error,
) *BridgeConnection {
	return &BridgeConnection{
		srcConn:      srcConn,
		dstConn:      dstConn,
		analysisChan: analysisChan,
		errChan:      errChan,
	}
}

// Bind starts the binding between the server and client connections.
func (b *BridgeConnection) Bind() {
	srcBuffer := bufio.NewReader(b.srcConn)
	dstBuffer := bufio.NewReader(b.dstConn)

	b.pipe(srcBuffer, dstBuffer)
}

func (b *BridgeConnection) pipe(src, dst *bufio.Reader) {
	go func() {
		for {
			b.copyWait(src, b.dstConn, b.errChan)
			b.copyWait(dst, b.srcConn, b.errChan)
		}
	}()
}

func (b *BridgeConnection) copyWait(src *bufio.Reader, dst net.Conn, errChan chan error) {
	data, err := b.read(src)

	if err != nil {
		errChan <- err
		errChan <- b.close()
		return
	}

	_, err = b.write(dst, []byte(data))

	if err != nil {
		errChan <- err
		errChan <- b.close()
		return
	}
}

func (b *BridgeConnection) read(src *bufio.Reader) ([]byte, error) {
	data, err := src.ReadString('\n')

	if err != nil {
		return []byte{}, err
	}

	return []byte(data), err
}

func (b *BridgeConnection) write(dst net.Conn, data []byte) (int, error) {
	// Updating analysis channel with the written data.
	b.analysisChan <- string(data)

	return dst.Write(data)
}

func (b *BridgeConnection) close() error {
	err := b.srcConn.Close()

	if err != nil {
		return err
	}

	err = b.dstConn.Close()

	return err
}
