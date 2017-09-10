package app

import (
	"bufio"
	"net"
	"unicode/utf8"
)

type (
	Connection struct {
		conn    net.Conn
		dataBus *DataBus
	}
)

func NewConnection(conn net.Conn, dataBus *DataBus) *Connection {
	return &Connection{
		conn:    conn,
		dataBus: dataBus,
	}
}

func (c *Connection) Start() {
	go c.read()
	go c.write()
}

func (c *Connection) read() {
	for data := range c.dataBus.ForwardChan {
		c.conn.Write([]byte(data))
		msgType := c.dataBus.GetMessageType(data)
		c.dataBus.Stats.UpdateTotalCounters(msgType)
		c.dataBus.LogTable.UpdateLog(msgType)
	}
}

func (c *Connection) write() {
	buffer := bufio.NewReader(c.conn)

	for {
		data, err := buffer.ReadString('\n')

		if err != nil {
			c.conn.Close()
			break
		}

		if utf8.RuneCountInString(data) > 0 {
			c.dataBus.ReverseChan <- data
			msgType := c.dataBus.GetMessageType(data)
			c.dataBus.Stats.UpdateTotalCounters(msgType)
			c.dataBus.LogTable.UpdateLog(msgType)
		}
	}
}
