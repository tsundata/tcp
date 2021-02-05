package main

import (
	"github.com/tsundata/tcp/network"
	"log"
)

type PingRouter struct {
	network.BaseRouter
}

func (r *PingRouter) BeforeHook(req network.IRequest) {
	log.Println("call router BeforeHook")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("before hook\n"))
	if err != nil {
		log.Println(err)
	}
}

func (r *PingRouter) Handle(req network.IRequest) {
	log.Println("call router Handle")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("handle\n"))
	if err != nil {
		log.Println(err)
	}
}

func (r *PingRouter) AfterHook(req network.IRequest) {
	log.Println("call router AfterHook")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("after hook\n"))
	if err != nil {
		log.Println(err)
	}
}

func main() {
	s := network.NewServer("example")
	s.AddRouter(&PingRouter{})
	s.Serve()
}
