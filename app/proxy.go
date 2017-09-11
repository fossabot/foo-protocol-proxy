package app

import (
	"fmt"
	"foo-protocol-proxy/config"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type (
	Proxy struct {
		config         config.Configuration
		clientConnChan chan net.Conn
		signalChan     chan os.Signal
		analyzer       *Analyzer
		milliTicker    *time.Ticker
		oneSecTicker   *time.Ticker
	}
)

func NewProxy(config config.Configuration) *Proxy {
	return &Proxy{
		config:         config,
		clientConnChan: make(chan net.Conn),
		signalChan:     make(chan os.Signal, 1),
		analyzer:       NewAnalyzer(make(AnalysisType)),
		milliTicker:    time.NewTicker(time.Millisecond),
		oneSecTicker:   time.NewTicker(time.Second),
	}
}

func (p *Proxy) Start() error {
	lis, err := net.Listen("tcp", p.config.Listening)

	if err != nil {
		return err
	}

	listener := NewListener(lis)

	log.Printf("Forwarding from %s to %s", listener.Addr(), p.config.Forwarding)

	go p.handleClientConnections(p.clientConnChan)
	go p.heartbeat()
	go p.reportStatus()
	go p.analyzer.monitorData()

	signal.Notify(p.signalChan, syscall.SIGUSR2)

	listener.awaitForConnections(p.clientConnChan)

	return nil
}

func (p *Proxy) handleClientConnections(clientConnChan chan net.Conn) {
	for clientConn := range clientConnChan {
		serverConn, err := p.Forward()

		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		bridgeConnection := NewBridgeConnection(clientConn, serverConn, p.analyzer.dataSrc)
		bridgeConnection.Bind()
	}

	p.Close()
}

func (p *Proxy) Forward() (net.Conn, error) {
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
			p.analyzer.timeTable.updateStats(time.Millisecond)

		case <-p.oneSecTicker.C:
			p.analyzer.timeTable.updateStats(time.Second)
		}
	}
}

func (p *Proxy) reportStatus() {
	for {
		<-p.signalChan
		p.analyzer.stats.CalculateAverages(p.analyzer.timeTable)
		report, err := p.analyzer.Report()

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(report)
	}
}

func (p *Proxy) Close() {
	close(p.clientConnChan)
	close(p.signalChan)
	p.milliTicker.Stop()
	p.oneSecTicker.Stop()
}
