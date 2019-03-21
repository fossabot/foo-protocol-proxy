package analysis

import (
	"sync"
	"sync/atomic"
	"time"
)

type (
	// TimeTable type holds reference for data that are updated over time.
	TimeTable struct {
		// Request time log for 1 second.
		RequestsInOneSec [1000]uint64
		// Response time log for 1 second.
		ResponsesInOneSec [1000]uint64
		// Request time log for 10 seconds.
		RequestsInTenSec [10]uint64
		// Response time log for 10 seconds.
		ResponsesInTenSec [10]uint64
		// Index for one second time log array.
		IndexOneSec uint16
		// Index for ten seconds time log array.
		IndexTenSec uint8
		// Holds total number of requests in one second.
		RequestsCount uint64
		// Holds total number of responses in one second.
		ResponsesCount uint64
	}
)

// UpdateCounters updates associated counter based on message type.
func (t *TimeTable) UpdateCounters(msgType message) {
	switch msgType {
	case TypeReq:
		atomic.AddUint64(&t.RequestsCount, 1)

	case TypeAck, TypeNak:
		atomic.AddUint64(&t.ResponsesCount, 1)
	}
}

func (t *TimeTable) updateStats(duration time.Duration) {
	mutex := sync.Mutex{}
	mutex.Lock()

	switch duration {
	case time.Millisecond:
		indexOneSec := t.IndexOneSec

		t.RequestsInOneSec[indexOneSec] = t.RequestsCount
		t.ResponsesInOneSec[indexOneSec] = t.ResponsesCount
		// Resetting requests counter.
		t.RequestsCount = 0
		// Resetting responses counter.
		t.ResponsesCount = 0
		// Updating the sliding window index.
		t.IndexOneSec++
		t.IndexOneSec %= 1000

	case time.Second:
		indexTenSec := t.IndexTenSec
		requestsSumOneSec, responsesSumOneSec := uint64(0), uint64(0)

		for _, val := range t.RequestsInOneSec {
			requestsSumOneSec += val
		}

		for _, val := range t.ResponsesInOneSec {
			responsesSumOneSec += val
		}

		t.RequestsInTenSec[indexTenSec] = requestsSumOneSec
		t.ResponsesInTenSec[indexTenSec] = responsesSumOneSec

		// Updating the sliding window index.
		t.IndexTenSec++
		t.IndexTenSec %= 10
	}

	mutex.Unlock()
}
