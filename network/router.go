package network

type IRouter interface {
	BeforeHook(IRequest)
	Handle(IRequest)
	AfterHook(IRequest)
}

type BaseRouter struct {
}

func (b *BaseRouter) BeforeHook(req IRequest) {
	panic("implement me")
}

func (b *BaseRouter) Handle(req IRequest) {
	panic("implement me")
}

func (b *BaseRouter) AfterHook(req IRequest) {
	panic("implement me")
}
