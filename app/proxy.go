package app

import (
	"foo-protocol-proxy/config"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type (
	Proxy struct {
		config    config.Configuration
		dataChan  chan *DataBus
		status    chan os.Signal
		stats     *Stats
		ticker    *time.Ticker
		timeTable *TimeTable
	}
)

func NewProxy(config config.Configuration) *Proxy {
	return &Proxy{
		config:   config,
		dataChan: make(chan *DataBus),
		status:   make(chan os.Signal, 1),
		stats:    new(Stats),
		ticker:   time.NewTicker(time.Millisecond),
		timeTable: &TimeTable{
			RequestsInOneSec:  [1000]uint64{},
			ResponseInOneSec:  [1000]uint64{},
			RequestsInTenSec:  [10000]uint64{},
			ResponsesInTenSec: [10000]uint64{},
		},
	}
}

func (p *Proxy) Start() {
	config := p.config
	listener, err := net.Listen("tcp", config.Listening)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Printf("Listening on %s, PID=%d", config.Listening, os.Getpid())

	go p.handleConnections()
	go p.heartbeat()

	signal.Notify(p.status, syscall.SIGUSR2)

	go p.reportStatus()

	for {
		clientConn, err := listener.Accept()

		if err != nil {
			log.Println(err)
			continue
		}

		dataBus := NewDataBus(p.stats, p.timeTable)
		NewClient(clientConn, dataBus).Start()

		p.dataChan <- dataBus
	}
}

func (p *Proxy) handleConnections() {
	for dataBus := range p.dataChan {
		serverConn := p.Forward()
		NewServer(serverConn, dataBus).Start()
	}
}

func (p *Proxy) Forward() net.Conn {
	serverConn, err := net.Dial("tcp", p.config.Forwarding)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return serverConn
}

func (p *Proxy) heartbeat() {
	for {
		select {
		case <-p.ticker.C:
			index := atomic.AddUint32(&p.timeTable.Index, 1) - 1
			// Index for one second time table.
			indexOneSec := index % 1000
			// Index for ten second time table.
			indexTenSec := index % 10000

			requestsCount := atomic.LoadUint64(&p.timeTable.RequestsCount)
			responsesCount := atomic.LoadUint64(&p.timeTable.ResponsesCount)

			mutex := sync.Mutex{}
			mutex.Lock()

			p.timeTable.RequestsInOneSec[indexOneSec] = requestsCount
			p.timeTable.ResponseInOneSec[indexOneSec] = responsesCount
			p.timeTable.RequestsInTenSec[indexTenSec] = requestsCount
			p.timeTable.ResponsesInTenSec[indexTenSec] = responsesCount

			p.timeTable.RequestsCount = 0
			p.timeTable.ResponsesCount = 0

			mutex.Unlock()
		}
	}
}

func (p *Proxy) reportStatus() {
	for {
		<-p.status
		p.stats.CalculateAverages(p.timeTable)
		p.stats.Flush()
	}
}
