package network

import (
	"fmt"
	"log"
	"net"
)

type ServerInterface interface {
	Start()
	Stop()
	Serve()
}

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
}

func NewServer(name string) *Server {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      5678,
	}
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
		for {
			conn, err := lis.AcceptTCP()
			if err != nil {
				log.Println(err)
				continue
			}
			// demo
			go func(conn net.Conn) {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						log.Println("recv buf err ", err)
						continue
					}
					log.Printf("recv data %s %d\n", buf, cnt)
					_, err = conn.Write(buf[:cnt])
					if err != nil {
						log.Println("write back buf err ", err)
						continue
					}
				}
			}(conn)
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
