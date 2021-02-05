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
	err := req.GetConnection().SendMessage(1, []byte("before hook\n"))
	if err != nil {
		log.Println(err)
	}
}

func (r *PingRouter) Handle(req network.IRequest) {
	log.Println("call router Handle")
	err := req.GetConnection().SendMessage(2, []byte("handle\n"))
	if err != nil {
		log.Println(err)
	}
}

func (r *PingRouter) AfterHook(req network.IRequest) {
	log.Println("call router AfterHook")
	err := req.GetConnection().SendMessage(3, []byte("after hook\n"))
	if err != nil {
		log.Println(err)
	}
}

func main() {
	s := network.NewServer("example")
	s.AddRouter(&PingRouter{})
	s.Serve()
}
