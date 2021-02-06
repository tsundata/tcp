package network

import (
	"errors"
	"fmt"
	"github.com/tsundata/tcp/utils"
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

	messageHandler IMessageHandle
	ExitBuffChan   chan bool
	messageChan    chan []byte
}

func NewConnection(conn *net.TCPConn, connID uint32, h IMessageHandle) *Connection {
	return &Connection{
		Conn:           conn,
		ConnID:         connID,
		isClosed:       false,
		messageHandler: h,
		ExitBuffChan:   make(chan bool, 1),
		messageChan:    make(chan []byte),
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
		if utils.Setting.WorkerPoolSize > 0 {
			c.messageHandler.SendMessageToTaskQueue(&req)
		} else {
			go c.messageHandler.DoMessageHandler(&req)
		}
	}
}

func (c *Connection) StartWriter() {
	log.Println("writer goroutine is running")
	defer log.Println(c.RemoteAddr(), " conn writer exit")

	for {
		select {
		case data := <-c.messageChan:
			if _, err := c.Conn.Write(data); err != nil {
				log.Println("send data error ", err, " conn writer exit")
				return
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}

func (c *Connection) Start() {
	go c.StartReader()
	go c.StartWriter()

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

	// write
	c.messageChan <- msg

	return nil
}
