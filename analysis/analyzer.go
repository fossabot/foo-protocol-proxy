package analysis

import (
	"encoding/json"
	"github.com/ahmedkamals/foo-protocol-proxy/persistance"
	"time"
)

type (
	Analyzer struct {
		dataSrc   AnalysisType
		stats     *Stats
		timeTable *TimeTable
	}

	AnalysisType chan string
	MessageType  uint
)

const (
	TYPE_REQ MessageType = iota
	TYPE_ACK
	TYPE_NAK
)

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		dataSrc: make(AnalysisType),
		stats:   new(Stats),
		timeTable: &TimeTable{
			RequestsInOneSec:  [1000]uint64{},
			ResponsesInOneSec: [1000]uint64{},
			RequestsInTenSec:  [10]uint64{},
			ResponsesInTenSec: [10]uint64{},
		},
	}
}

func (a *Analyzer) RestoreTenSecCounter(recovery *persistance.Recovery) {
	if recovery != nil {
		diff := uint64(time.Now().Unix()) - recovery.TimeStamp

		if diff <= 10 {
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
	}
}

func (a *Analyzer) GetMessageType(msg string) MessageType {
	msgPrefix := msg[0:3]
	msgType := TYPE_REQ

	switch msgPrefix {
	case "REQ":
		msgType = TYPE_REQ

	case "ACK":
		msgType = TYPE_ACK

	case "NAK":
		msgType = TYPE_NAK
	}

	return msgType
}

func (a *Analyzer) MonitorData() {
	for data := range a.dataSrc {
		msgType := a.GetMessageType(data)
		a.stats.UpdateTotalCounters(msgType)
		a.timeTable.UpdateCounters(msgType)
	}
}

func (a *Analyzer) Report() (string, error) {
	a.calculateAverages()
	result, err := json.Marshal(a.stats)

	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (a *Analyzer) GetDataSource() AnalysisType {
	return a.dataSrc
}

func (a *Analyzer) GetTimeTable() *TimeTable {
	return a.timeTable
}

func (a *Analyzer) UpdateStats(duration time.Duration) {
	a.timeTable.updateStats(duration)
}

func (a *Analyzer) calculateAverages() {
	a.stats.CalculateAverages(a.timeTable)
}
