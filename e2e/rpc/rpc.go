package rpc

import "context"

type RpcGreetingServiceImpl struct {
}

var _ RpcGreetingService = (*RpcGreetingServiceImpl)(nil)

func (s *RpcGreetingServiceImpl) SayHello(ctx context.Context, name string) (string, error) {
	return "Hello " + name, nil
}
