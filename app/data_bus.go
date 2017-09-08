package app

type (
	DataBus struct {
		ForwardChan chan string
		ReverseChan chan string
		Stats       *Stats
		LogTable    *TimeTable
	}

	MessageType uint
)

const (
	TYPE_REQ MessageType = iota
	TYPE_ACK
	TYPE_NAK
)

func NewDataBus(stats *Stats, logTable *TimeTable) *DataBus {
	return &DataBus{
		ForwardChan: make(chan string),
		ReverseChan: make(chan string),
		Stats:       stats,
		LogTable:    logTable,
	}
}

func (d *DataBus) GetMessageType(msg string) MessageType {
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
