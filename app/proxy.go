package app

import (
	"github.com/ahmedkamals/foo-protocol-proxy/analysis"
	"github.com/ahmedkamals/foo-protocol-proxy/communication"
	"github.com/ahmedkamals/foo-protocol-proxy/config"
	"github.com/ahmedkamals/foo-protocol-proxy/persistence"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type (
	// Proxy orchestrates the interactions between the server and client,
	// and collection of analysis data.
	Proxy struct {
		config         config.Configuration
		clientConnChan chan net.Conn
		analyzer       *analysis.Analyzer
		signalChan     chan os.Signal
		errorChan      chan error
		milliTicker    *time.Ticker
		oneSecTicker   *time.Ticker
		saver          *persistence.Saver
	}
)

// NewProxy allocates and returns a new Proxy to handle connections forwarding/reversing.
func NewProxy(
	config config.Configuration,
	analyzer *analysis.Analyzer,
	saver *persistence.Saver,
) *Proxy {
	return &Proxy{
		config:         config,
		clientConnChan: make(chan net.Conn),
		analyzer:       analyzer,
		signalChan:     make(chan os.Signal, 1),
		errorChan:      make(chan error, 10),
		milliTicker:    time.NewTicker(time.Millisecond),
		oneSecTicker:   time.NewTicker(time.Second),
		saver:          saver,
	}
}

// Start initiates the proxy operations.
func (p *Proxy) Start() error {
	lis, err := net.Listen("tcp", p.config.Listening)

	if err != nil {
		return err
	}

	p.recoverData()

	listener := communication.NewListener(lis, p.errorChan)

	log.Printf("Forwarding from %s to %s", listener.Addr(), p.config.Forwarding)

	go listener.AwaitForConnections(p.clientConnChan)
	go p.handleClientConnections(p.clientConnChan)
	go p.heartbeat()
	go p.reportStatus()
	go p.analyzer.MonitorData()
	go p.monitorErrors()

	signal.Notify(p.signalChan, syscall.SIGUSR2)

	return nil
}

func (p *Proxy) recoverData() {
	data, err := p.saver.Read()

	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	recovery := persistence.NewEmptyRecovery()
	recovery.Unmarshal(data)

	mutex := sync.Mutex{}
	mutex.Lock()
	p.analyzer.RestoreTenSecCounter(recovery)
	mutex.Unlock()
}

func (p *Proxy) handleClientConnections(clientConnChan chan net.Conn) {
	for clientConn := range clientConnChan {
		serverConn, err := p.forward()

		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		bridgeConnection := communication.NewBridgeConnection(
			clientConn,
			serverConn,
			p.analyzer.GetDataSource(),
			p.errorChan,
		)
		bridgeConnection.Bind()
	}
}

func (p *Proxy) forward() (net.Conn, error) {
	serverConn, err := net.Dial("tcp", p.config.Forwarding)

	if err != nil {
		return nil, err
	}

	return serverConn, nil
}

func (p *Proxy) heartbeat() {
	for {
		select {
		case <-p.milliTicker.C:
			p.analyzer.UpdateStats(time.Millisecond)

		case <-p.oneSecTicker.C:
			p.analyzer.UpdateStats(time.Second)
			p.persistData()
		}
	}
}

func (p *Proxy) persistData() {
	timeTable := p.analyzer.GetTimeTable()
	r := persistence.NewRecovery(
		timeTable.IndexTenSec,
		uint64(time.Now().Unix()),
		timeTable.RequestsInTenSec,
		timeTable.ResponsesInTenSec,
	)
	data, err := r.Marshal()

	if err != nil {
		log.Fatal(err)
	}

	p.saver.Save(data)
}

func (p *Proxy) reportStatus() {
	for {
		<-p.signalChan
		report, err := p.analyzer.Report()

		if err != nil {
			p.errorChan <- err
			return
		}

		println(report)
	}
}

func (p *Proxy) monitorErrors() {
	for {
		select {
		case err := <-p.errorChan:
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (p *Proxy) close() {
	close(p.clientConnChan)
	close(p.signalChan)
	p.saver.Close()
	p.milliTicker.Stop()
	p.oneSecTicker.Stop()
}
