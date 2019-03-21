package analysis

import (
	"encoding/json"
	"github.com/ahmedkamals/foo-protocol-proxy/app/persistence"
	"time"
)

type (
	// Analyzer will hold analysis data.
	Analyzer struct {
		dataSrc   Aggregation
		stats     *Stats
		timeTable *TimeTable
	}
	// Aggregation used for intercepted data.
	Aggregation chan string
	// message used for differentiating message types.
	message uint
)

const (
	// TypeReq stands for "REQ"
	TypeReq message = iota
	// TypeAck stands for "ACK"
	TypeAck
	// TypeNak stands for "NAK"
	TypeNak
)

// NewAnalyzer allocates and returns a new Analyzer, that intercepts data, and aggregates it.
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		dataSrc: make(Aggregation),
		stats:   new(Stats),
		timeTable: &TimeTable{
			RequestsInOneSec:  [1000]uint64{},
			ResponsesInOneSec: [1000]uint64{},
			RequestsInTenSec:  [10]uint64{},
			ResponsesInTenSec: [10]uint64{},
		},
	}
}

// RestoreTenSecCounter restores the value of the 10s second counter from recovery data.
// Considering the time that has passed since the last failure.
func (a *Analyzer) RestoreTenSecCounter(recovery *persistence.Recovery) {
	if recovery != nil {
		return
	}

	diff := uint64(time.Now().Unix()) - recovery.TimeStamp
	if diff > 10 {
		return
	}

	a.timeTable.IndexTenSec = recovery.Index
	// Restore only the parts we are interested in.
	rangeStart := uint64(a.timeTable.IndexTenSec)
	rangeEnd := (diff + rangeStart) % 10

	for i := uint64(0); i < 10; i++ {
		// Check if overlapping period
		// [##111001##]
		// # means a cancelled value, as the proxy stopped at second 9
		// And comes back to work at second 3.
		if (rangeEnd < rangeStart &&
			(i >= rangeStart || i <= rangeEnd)) ||
			i >= rangeStart && i <= rangeEnd {
			continue
		}

		a.timeTable.RequestsInTenSec[i] = recovery.RequestsInTenSec[i]
		a.timeTable.ResponsesInTenSec[i] = recovery.ResponsesInTenSec[i]
	}
}

func (a *Analyzer) getMessageType(msg string) message {
	msgPrefix := msg[0:3]
	msgType := TypeReq

	switch msgPrefix {
	case "REQ":
		msgType = TypeReq

	case "ACK":
		msgType = TypeAck

	case "NAK":
		msgType = TypeNak
	}

	return msgType
}

// MonitorData initiates data monitoring to be aggregated.
func (a *Analyzer) MonitorData() {
	for data := range a.dataSrc {
		msgType := a.getMessageType(data)
		a.stats.UpdateTotalCounters(msgType)
		a.timeTable.UpdateCounters(msgType)
	}
}

// Report reports the aggregated data as json string.
func (a *Analyzer) Report() (string, error) {
	a.calculateAverages()
	result, err := json.Marshal(a.stats)

	if err != nil {
		return "", err
	}

	return string(result), nil
}

// GetDataSource returns the aggregation chanel.
func (a *Analyzer) GetDataSource() Aggregation {
	return a.dataSrc
}

// GetTimeTable returns a reference to a TimeTable instance.
func (a *Analyzer) GetTimeTable() *TimeTable {
	return a.timeTable
}

// UpdateStats acts as a wrapper for updating TimeTable.
func (a *Analyzer) UpdateStats(duration time.Duration) {
	a.timeTable.updateStats(duration)
}

func (a *Analyzer) calculateAverages() {
	a.stats.CalculateAverages(a.timeTable)
}
