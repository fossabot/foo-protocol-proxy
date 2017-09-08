package app

import (
	"bufio"
	"net"
	"unicode/utf8"
)

type (
	Client struct {
		conn    net.Conn
		dataBus *DataBus
	}
)

func NewClient(conn net.Conn, dataBus *DataBus) *Client {
	return &Client{
		conn:    conn,
		dataBus: dataBus,
	}
}

func (c *Client) Start() {
	go c.read()
	go c.write()
}

func (c *Client) read() {
	for data := range c.dataBus.ReverseChan {
		c.conn.Write([]byte(data))
	}
}

func (c *Client) write() {
	buffer := bufio.NewReader(c.conn)

	for {
		str, err := buffer.ReadString('\n')

		if err != nil {
			c.conn.Close()
			break
		}

		if utf8.RuneCountInString(str) > 0 {
			c.dataBus.ForwardChan <- str
		}
	}
}
