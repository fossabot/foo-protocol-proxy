package app

import (
	"sync/atomic"
)

type (
	TimeTable struct {
		// Request time log for 1 second.
		RequestsInOneSec [1000]uint64
		// Response time log for 1 second.
		ResponseInOneSec [1000]uint64
		// Request time log for 10 seconds.
		RequestsInTenSec [10000]uint64
		// Response time log for 10 seconds.
		ResponsesInTenSec [10000]uint64
		RequestsCount     uint64
		ResponsesCount    uint64
		Index             uint32
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
