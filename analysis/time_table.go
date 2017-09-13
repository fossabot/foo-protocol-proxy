package analysis

import (
	"sync"
	"sync/atomic"
	"time"
)

type (
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
		IndexOneSec uint32
		// Index for ten seconds time log array.
		IndexTenSec uint32
		// Holds total number of requests in one second.
		RequestsCount uint64
		// Holds total number of responses in one second.
		ResponsesCount uint64
	}
)

func (l *TimeTable) UpdateCounters(msgType MessageType) {
	switch msgType {
	case TYPE_REQ:
		atomic.AddUint64(&l.RequestsCount, 1)

	case TYPE_ACK, TYPE_NAK:
		atomic.AddUint64(&l.ResponsesCount, 1)
	}
}

func (l *TimeTable) updateStats(duration time.Duration) {
	mutex := sync.Mutex{}
	mutex.Lock()

	switch duration {
	case time.Millisecond:
		indexOneSec := l.IndexOneSec

		l.RequestsInOneSec[indexOneSec] = l.RequestsCount
		l.ResponsesInOneSec[indexOneSec] = l.ResponsesCount
		// Resetting requests counter.
		l.RequestsCount = 0
		// Resetting responses counter.
		l.ResponsesCount = 0
		// Updating the sliding window index.
		l.IndexOneSec++
		l.IndexOneSec %= 1000

	case time.Second:
		indexTenSec := l.IndexTenSec
		requestsSumOneSec, responsesSumOneSec := uint64(0), uint64(0)

		for _, val := range l.RequestsInOneSec {
			requestsSumOneSec += val
		}

		for _, val := range l.ResponsesInOneSec {
			responsesSumOneSec += val
		}

		l.RequestsInTenSec[indexTenSec] = requestsSumOneSec
		l.ResponsesInTenSec[indexTenSec] = responsesSumOneSec

		// Updating the sliding window index.
		l.IndexTenSec++
		l.IndexTenSec %= 10
	}

	mutex.Unlock()
}
