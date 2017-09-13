package analysis

import (
	"encoding/json"
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

func (l *Analyzer) GetMessageType(msg string) MessageType {
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

func (l *Analyzer) MonitorData() {
	for data := range l.dataSrc {
		msgType := l.GetMessageType(data)
		l.stats.UpdateTotalCounters(msgType)
		l.timeTable.UpdateCounters(msgType)
	}
}

func (l *Analyzer) Report() (string, error) {
	l.calculateAverages()
	result, err := json.Marshal(l.stats)

	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (l *Analyzer) GetDataSource() AnalysisType {
	return l.dataSrc
}

func (l *Analyzer) UpdateStats(duration time.Duration) {
	l.timeTable.updateStats(duration)
}

func (l *Analyzer) calculateAverages() {
	l.stats.CalculateAverages(l.timeTable)
}
