package app

import (
	"encoding/json"
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

func NewAnalyzer(AnalysisChan AnalysisType) *Analyzer {
	return &Analyzer{
		dataSrc: AnalysisChan,
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

func (l *Analyzer) monitorData() {
	for data := range l.dataSrc {
		msgType := l.GetMessageType(data)
		l.stats.UpdateTotalCounters(msgType)
		l.timeTable.UpdateCounters(msgType)
	}
}

func (l *Analyzer) Report() (string, error) {
	result, err := json.Marshal(l.stats)

	if err != nil {
		return "", err
	}

	return string(result), nil
}
