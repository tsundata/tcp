package network

import (
	"fmt"
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
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      5678,
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
	log.Printf("[TS] server starting... %s:%d", s.IP, s.Port)
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
