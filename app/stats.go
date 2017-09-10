package app

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
)

type (
	Stats struct {
		// msg_total — total number of messages (integer)
		TotalMessages uint64 `json:"msg_total"`
		// msg_req — total number of REQ messages (integer)
		TotalRequests uint64 `json:"msg_req"`
		// msg_ack — total number of ACK messages (integer)
		TotalACK uint64 `json:"msg_ack"`
		// msg_nak — total number of NAK messages (integer)
		TotalNAK uint64 `json:"msg_nak"`
		// request_rate_1s — average REQ messages/sec, in a 1s moving window (floating point)
		RequestRatePerSecond float64 `json:"request_rate_1s"`
		// request_rate_10s — average REQ messages/sec, in a 10s moving window (floating point)
		RequestRatePerTenSecond float64 `json:"request_rate_10s"`
		// response_rate_1s — average ACK+NAK messages per second, in a 1s moving window (floating point)
		ResponseRatePerSecond float64 `json:"response_rate_1s"`
		// response_rate_10s — average ACK+NAK messages per second, in a 10s moving window (floating point)
		ResponseRatePerTenSecond float64 `json:"response_rate_10s"`
	}
)

func (s *Stats) UpdateTotalCounters(msgType MessageType) {
	switch msgType {
	case TYPE_REQ:
		atomic.AddUint64(&s.TotalRequests, 1)

	case TYPE_ACK:
		atomic.AddUint64(&s.TotalACK, 1)

	case TYPE_NAK:
		atomic.AddUint64(&s.TotalNAK, 1)
	}

	atomic.AddUint64(&s.TotalMessages, 1)
}

// Calculates average req/response in 1s/10s.
func (s *Stats) CalculateAverages(timeTable *TimeTable) {
	mutex := sync.Mutex{}
	mutex.Lock()

	requestsSumOneSec := uint64(0)
	responsesSumOneSec := uint64(0)
	requestsSumTenSec := uint64(0)
	responsesSumTenSec := uint64(0)

	for _, val := range timeTable.RequestsInOneSec {
		requestsSumOneSec += val
	}

	for _, val := range timeTable.ResponsesInOneSec {
		responsesSumOneSec += val
	}

	for _, val := range timeTable.RequestsInTenSec {
		requestsSumTenSec += val
	}

	for _, val := range timeTable.ResponsesInTenSec {
		responsesSumTenSec += val
	}

	s.RequestRatePerSecond = float64(requestsSumOneSec) / 1000
	s.ResponseRatePerSecond = float64(responsesSumOneSec) / 1000
	s.RequestRatePerTenSecond = float64(requestsSumTenSec) / 10000
	s.ResponseRatePerTenSecond = float64(responsesSumTenSec) / 10000

	mutex.Unlock()
}

func (s *Stats) Flush() {
	result, err := json.Marshal(s)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(result))
}
