package network

import (
	"fmt"
	"github.com/tsundata/tcp/utils"
	"log"
	"net"
)

type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(uint32, IRouter)
	GetConnManager() IConnManager
	SetOnConnStart(func(IConnection))
	SetOnConnStop(func(IConnection))
	CallOnConnStart(IConnection)
	CallOnConnStop(IConnection)
}

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int

	messageHandler IMessageHandle
	ConnManager    IConnManager

	OnConnStart func(conn IConnection)
	OnConnStop  func(conn IConnection)
}

func NewServer() IServer {
	utils.Setting.Reload()
	return &Server{
		Name:           utils.Setting.Name,
		IPVersion:      "tcp4",
		IP:             utils.Setting.Host,
		Port:           utils.Setting.TCPPort,
		messageHandler: NewMessageHandle(),
		ConnManager:    NewConnManager(),
	}
}

func (s *Server) Start() {
	log.Printf("[TS] server starting...%s %s:%d", s.Name, s.IP, s.Port)
	log.Printf("[TS] version: %s maxConn: %d maxPacketSize: %d",
		utils.Setting.Version,
		utils.Setting.MaxConn,
		utils.Setting.MaxPacketSize,
	)
	go func() {
		// start worker pool
		s.messageHandler.StartWorkerPool()

		// addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			log.Println(err)
			return
		}
		// listen
		lis, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			log.Println(err)
			return
		}

		// accept
		var cid uint32
		cid = 0
		for {
			conn, err := lis.AcceptTCP()
			if err != nil {
				log.Println(err)
				continue
			}

			if s.ConnManager.Len() >= utils.Setting.MaxConn {
				_ = conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, cid, s.messageHandler)
			cid++

			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	log.Println("[TS] stop server, name ", s.Name)
	s.ConnManager.ClearConn()
}

func (s *Server) Serve() {
	s.Start()

	select {}
}

func (s *Server) AddRouter(id uint32, r IRouter) {
	s.messageHandler.AddRouter(id, r)
	log.Println("add router success")
}

func (s *Server) GetConnManager() IConnManager {
	return s.ConnManager
}

func (s *Server) SetOnConnStart(f func(IConnection)) {
	s.OnConnStart = f
}

func (s *Server) SetOnConnStop(f func(IConnection)) {
	s.OnConnStop = f
}

func (s *Server) CallOnConnStart(conn IConnection) {
	if s.OnConnStart != nil {
		log.Println("call on conn start...")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn IConnection) {
	if s.OnConnStop != nil {
		log.Println("call on conn stop...")
		s.OnConnStop(conn)
	}
}
