package network

import (
	"log"
	"strconv"
)

type IMessageHandle interface {
	DoMessageHandler(IRequest)
	AddRouter(uint32, IRouter)
}

type MessageHandle struct {
	APIs map[uint32]IRouter
}

func NewMessageHandle() *MessageHandle {
	return &MessageHandle{
		APIs: make(map[uint32]IRouter),
	}
}

func (m MessageHandle) DoMessageHandler(req IRequest) {
	handler, ok := m.APIs[req.GetMessageID()]
	if !ok {
		log.Println("api message id = ", req.GetMessageID(), " is not found")
		return
	}
	handler.BeforeHook(req)
	handler.Handle(req)
	handler.AfterHook(req)
}

func (m MessageHandle) AddRouter(id uint32, router IRouter) {
	if _, ok := m.APIs[id]; ok {
		panic("repeated api, id = " + strconv.Itoa(int(id)))
	}
	m.APIs[id] = router
	log.Println("add api id = ", id)
}
