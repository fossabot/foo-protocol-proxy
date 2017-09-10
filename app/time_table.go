package app

import (
	"sync/atomic"
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
		RequestsCount     uint64
		// Holds total number of responses in one second.
		ResponsesCount    uint64
	}
)

func (l *TimeTable) UpdateLog(msgType MessageType) {
	switch msgType {
	case TYPE_REQ:
		atomic.AddUint64(&l.RequestsCount, 1)

	case TYPE_ACK, TYPE_NAK:
		atomic.AddUint64(&l.ResponsesCount, 1)
	}
}
