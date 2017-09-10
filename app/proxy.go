package app

import (
	"foo-protocol-proxy/config"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type (
	Proxy struct {
		config       config.Configuration
		dataChan     chan *DataBus
		status       chan os.Signal
		stats        *Stats
		milliTicker  *time.Ticker
		oneSecTicker *time.Ticker
		timeTable    *TimeTable
	}
)

func NewProxy(config config.Configuration) *Proxy {
	return &Proxy{
		config:       config,
		dataChan:     make(chan *DataBus),
		status:       make(chan os.Signal, 1),
		stats:        new(Stats),
		milliTicker:  time.NewTicker(time.Millisecond),
		oneSecTicker: time.NewTicker(time.Second),
		timeTable: &TimeTable{
			RequestsInOneSec:  [1000]uint64{},
			ResponsesInOneSec: [1000]uint64{},
			RequestsInTenSec:  [10]uint64{},
			ResponsesInTenSec: [10]uint64{},
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

	log.Printf("Forwarding from %s to %s", config.Listening, config.Forwarding)

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
		NewConnection(serverConn, dataBus).Start()
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
		case <-p.milliTicker.C:
			p.updateStats(time.Millisecond)

		case <- p.oneSecTicker.C:
			p.updateStats(time.Second)
		}
	}
}

func (p *Proxy) updateStats(duration time.Duration)  {
	mutex := sync.Mutex{}
	mutex.Lock()

	switch duration {
	case time.Millisecond:
		indexOneSec := p.timeTable.IndexOneSec

		p.timeTable.RequestsInOneSec[indexOneSec] = p.timeTable.RequestsCount
		p.timeTable.ResponsesInOneSec[indexOneSec] = p.timeTable.ResponsesCount
		// Resetting requests counter.
		p.timeTable.RequestsCount = 0
		// Resetting responses counter.
		p.timeTable.ResponsesCount = 0
		// Updating the sliding window index.
		p.timeTable.IndexOneSec++
		p.timeTable.IndexOneSec %= 1000

	case time.Second:
		indexTenSec := p.timeTable.IndexTenSec
		requestsSumOneSec, responsesSumOneSec := uint64(0), uint64(0)

		for _, val := range p.timeTable.RequestsInOneSec {
			requestsSumOneSec += val
		}

		for _, val := range p.timeTable.ResponsesInOneSec {
			responsesSumOneSec += val
		}

		p.timeTable.RequestsInTenSec[indexTenSec] = requestsSumOneSec
		p.timeTable.ResponsesInTenSec[indexTenSec] = responsesSumOneSec

		// Updating the sliding window index.
		p.timeTable.IndexTenSec++
		p.timeTable.IndexTenSec %= 10
	}

	mutex.Unlock()
}

func (p *Proxy) reportStatus() {
	for {
		<-p.status
		p.stats.CalculateAverages(p.timeTable)
		p.stats.Flush()
	}
}
