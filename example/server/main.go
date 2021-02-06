package main

import (
	"fmt"
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

func DoConnBegin(conn network.IConnection) {
	log.Println("DoConnBegin ...")

	log.Println("set conn name")
	conn.SetProperty("name", "demo")

	err := conn.SendMessage(2, []byte("DoConnBegin..."))
	if err != nil {
		log.Println(err)
	}
}

func DoConnLost(conn network.IConnection) {
	log.Println("DoConnLost ...")
	if name, err := conn.GetProperty("name"); err == nil {
		fmt.Println("conn property name = ", name)
	}
}

func main() {
	s := network.NewServer()

	s.SetOnConnStart(DoConnBegin)
	s.SetOnConnStop(DoConnLost)

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	s.Serve()
}
