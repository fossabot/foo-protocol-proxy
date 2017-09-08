package app

import (
	"bufio"
	"net"
	"unicode/utf8"
)

type (
	Server struct {
		conn    net.Conn
		dataBus *DataBus
	}
)

func NewServer(conn net.Conn, dataBus *DataBus) *Server {
	return &Server{
		conn:    conn,
		dataBus: dataBus,
	}
}

func (s *Server) Start() {
	go s.read()
	go s.write()
}

func (s *Server) read() {
	for data := range s.dataBus.ForwardChan {
		s.conn.Write([]byte(data))
		msgType := s.dataBus.GetMessageType(data)
		s.dataBus.Stats.UpdateTotalCounters(msgType)
		s.dataBus.LogTable.UpdateLog(msgType)
	}
}

func (s *Server) write() {
	buffer := bufio.NewReader(s.conn)

	for {
		data, err := buffer.ReadString('\n')

		if err != nil {
			break
		}

		if utf8.RuneCountInString(data) > 0 {
			s.dataBus.ReverseChan <- data
			msgType := s.dataBus.GetMessageType(data)
			s.dataBus.Stats.UpdateTotalCounters(msgType)
			s.dataBus.LogTable.UpdateLog(msgType)
		}
	}
}
