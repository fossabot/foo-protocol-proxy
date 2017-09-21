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
	// Reading buffer for source
	bufSrcReader *bufio.Reader
	// Reading buffer for destination
	bufDstReader *bufio.Reader
	// Analysis channel
	analysisChan analysis.AnalysisType
}

func NewBridgeConnection(srcConn, dstConn net.Conn, analysisChan analysis.AnalysisType) *BridgeConnection {
	return &BridgeConnection{
		srcConn:      srcConn,
		dstConn:      dstConn,
		bufSrcReader: bufio.NewReader(srcConn),
		bufDstReader: bufio.NewReader(dstConn),
		analysisChan: analysisChan,
	}
}

func (b *BridgeConnection) Bind() {
	go func() {
		for {
			b.Forward()
			b.Reverse()
		}
	}()
}

func (b *BridgeConnection) Forward() error {
	data, err := b.read(b.bufSrcReader)

	if err != nil {
		b.Close()
		return err
	}

	_, err = b.write(b.dstConn, []byte(data))

	if err != nil {
		b.Close()
		return err
	}

	return nil
}

func (b *BridgeConnection) Reverse() error {
	data, err := b.read(b.bufDstReader)

	if err != nil {
		b.Close()
		return err
	}

	_, err = b.write(b.srcConn, []byte(data))

	if err != nil {
		b.Close()
		return err
	}

	return nil
}

func (b *BridgeConnection) read(buffer *bufio.Reader) (string, error) {
	str, err := buffer.ReadString('\n')

	if err != nil {
		return "", err
	}

	return str, err
}

func (b *BridgeConnection) write(conn net.Conn, data []byte) (int, error) {
	// Updating analysis channel with the written data.
	b.analysisChan <- string(data)

	return conn.Write(data)
}

func (b *BridgeConnection) Close() error {
	err := b.srcConn.Close()

	if err != nil {
		return err
	}

	err = b.dstConn.Close()

	return err
}
