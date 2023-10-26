package main

import (
	"context"
	service "kitex_study/kitex_gen/kitex_gen/service"
)

// HelloImpl implements the last service interface defined in the IDL.
type HelloImpl struct{}

// Send implements the HelloImpl interface.
func (s *HelloImpl) Send(ctx context.Context, req *service.HelloReq) (resp *service.HelloResp, err error) {
	resp = &service.HelloResp{Res: "hello " + req.Name}
	return
}
