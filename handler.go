package main

import (
	"context"
	"kitex_study/kitex_gen/kitex_gen/service"
	"time"
)

// HelloImpl implements the last service interface defined in the IDL.
type HelloImpl struct{}

// Send implements the HelloImpl interface.
func (s *HelloImpl) Send(ctx context.Context, req *service.HelloReq) (resp *service.HelloResp, err error) {
	time.Sleep(time.Second * 3)
	resp = &service.HelloResp{Res: "hello " + req.Name}
	return
}
