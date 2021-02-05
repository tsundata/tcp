package network

import (
	"fmt"
	"log"
	"net"
)

type IConnection interface {
	Start()
	Stop()
	GetTCPConnection() *net.TCPConn
	GetConnID() uint32
	RemoteAddr() net.Addr
}

type HandFunc func(*net.TCPConn, []byte, int) error

type Connection struct {
	Conn     *net.TCPConn
	ConnID   uint32
	isClosed bool

	Router IRouter

	ExitBuffChan chan bool
}

func NewConnection(conn *net.TCPConn, connID uint32, r IRouter) *Connection {
	return &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		Router:       r,
		ExitBuffChan: make(chan bool, 1),
	}
}

func (c *Connection) StartReader() {
	log.Printf("reader goroutine is running... %d", c.ConnID)
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit")
	defer c.Stop()

	for {
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			log.Println("recv buf err ", err)
			c.ExitBuffChan <- true
			continue
		}

		req := Request{
			conn: c,
			data: buf,
		}
		go func(req IRequest) {
			c.Router.BeforeHook(req)
			c.Router.Handle(req)
			c.Router.AfterHook(req)
		}(&req)
	}
}

func (c *Connection) Start() {
	go c.StartReader()

	for {
		select {
		case <-c.ExitBuffChan:
			return
		}
	}
}

func (c *Connection) Stop() {
	if c.isClosed {
		return
	}
	c.isClosed = true

	err := c.Conn.Close()
	if err != nil {
		log.Println(err)
		return
	}

	c.ExitBuffChan <- true

	close(c.ExitBuffChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}