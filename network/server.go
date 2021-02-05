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
	AddRouter(router IRouter)
}

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int

	Router IRouter
}

func NewServer(name string) IServer {
	utils.Setting.Reload()
	return &Server{
		Name:      utils.Setting.Name,
		IPVersion: "tcp4",
		IP:        utils.Setting.Host,
		Port:      utils.Setting.TCPPort,
	}
}

func CallbackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	log.Printf("[TS] conn handle...")
	_, err := conn.Write(data[:cnt])
	if err != nil {
		log.Println("write back buf err ", err)
		return err
	}
	return nil
}

func (s *Server) Start() {
	log.Printf("[TS] server starting...%s %s:%d", s.Name, s.IP, s.Port)
	log.Printf("[TS] version: %s maxConn: %d maxPacketSize: %d",
		utils.Setting.Version,
		utils.Setting.MaxConn,
		utils.Setting.MaxPacketSize,
	)
	go func() {
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			log.Println(err)
			return
		}
		lis, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			log.Println(err)
			return
		}

		var cid uint32
		cid = 0

		for {
			conn, err := lis.AcceptTCP()
			if err != nil {
				log.Println(err)
				continue
			}

			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	// TODO
}

func (s *Server) Serve() {
	s.Start()

	select {}
}

func (s *Server) AddRouter(r IRouter) {
	s.Router = r
	log.Println("add router success")
}
