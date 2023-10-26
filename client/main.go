package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/remote/trans/nphttp2/status"
	"github.com/cloudwego/kitex/transport"
	"kitex_study/kitex_gen/kitex_gen/service"
	"kitex_study/kitex_gen/kitex_gen/service/hello"
)

func main() {
	newClient, err := hello.NewClient("hello", client.WithHostPorts("0.0.0.0:8888"), client.WithTransportProtocol(transport.GRPC))
	if err != nil {
		fmt.Println(err)
	}
	req := &service.HelloReq{Name: "zhr"}
	send, err := newClient.Send(context.Background(), req)
	if err != nil {
		if bizErr, ok := kerrors.FromBizStatusError(err); ok {
			fmt.Println(bizErr.BizStatusCode())
			fmt.Println(bizErr.BizMessage())
			println(bizErr.(status.Iface).GRPCStatus().Details()[0].(*service.HelloResp).Res)
		}
	}
	fmt.Println(send)
}
