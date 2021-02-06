package main

import (
	"github.com/tsundata/tcp/network"
	"log"
)

type PingRouter struct {
	network.BaseRouter
}

func (r *PingRouter) Handle(req network.IRequest) {
	log.Println("call router Handle")
	err := req.GetConnection().SendMessage(1, []byte("ping\n"))
	if err != nil {
		log.Println(err)
	}
}

type HelloRouter struct {
	network.BaseRouter
}

func (r *HelloRouter) Handle(req network.IRequest) {
	log.Println("call router Handle")
	err := req.GetConnection().SendMessage(2, []byte("hello\n"))
	if err != nil {
		log.Println(err)
	}
}

func DoConnBegin(conn network.IConnection)  {
	log.Println("DoConnBegin ...")
	err := conn.SendMessage(2, []byte("DoConnBegin..."))
	if err != nil {
		log.Println(err)
	}
}

func DoConnLost(conn network.IConnection)  {
	log.Println("DoConnLost ...")
}

func main() {
	s := network.NewServer()

	s.SetOnConnStart(DoConnBegin)
	s.SetOnConnStop(DoConnLost)

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	s.Serve()
}
