package network

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type IConnection interface {
	Start()
	Stop()
	GetTCPConnection() *net.TCPConn
	GetConnID() uint32
	RemoteAddr() net.Addr
	SendMessage(uint32, []byte) error
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
		pack := NewPack()

		headData := make([]byte, pack.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			log.Println("read message head error ", err)
			c.ExitBuffChan <- true
			continue
		}

		msg, err := pack.Unpack(headData)
		if err != nil {
			log.Println("unpack error ", err)
			c.ExitBuffChan <- true
			continue
		}

		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				log.Println("read message data error ", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		msg.SetData(data)

		req := Request{
			conn: c,
			data: msg,
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

func (c *Connection) SendMessage(id uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when send message")
	}

	pack := NewPack()
	msg, err := pack.Pack(NewMessage(id, data))
	if err != nil {
		log.Println("pack error message id = ", id)
		return errors.New("pack error message")
	}

	if _, err := c.Conn.Write(msg); err != nil {
		log.Println("write message id ", id, " error")
		c.ExitBuffChan <- true
		return errors.New("conn write error")
	}
	return nil
}
