package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/remote/trans/nphttp2/status"
	etcd "github.com/kitex-contrib/registry-etcd"
	"kitex_study/kitex_gen/kitex_gen/service"
	"kitex_study/kitex_gen/kitex_gen/service/hello"
	"time"
)

func main() {
	registry, err := etcd.NewEtcdResolver([]string{"127.0.0.1:12380"})
	if err != nil {
		panic(err)
	}
	newClient, err := hello.NewClient("Hello", client.WithResolver(registry))
	if err != nil {
		fmt.Println(err)
	}
	req := &service.HelloReq{Name: "jjking"}
	send, err := newClient.Send(context.Background(), req, callopt.WithRPCTimeout(time.Second*2))
	if err != nil {
		if bizErr, ok := kerrors.FromBizStatusError(err); ok {
			fmt.Println(bizErr.BizStatusCode())
			fmt.Println(bizErr.BizMessage())
			fmt.Println(bizErr.(status.Iface).GRPCStatus().Details()[0].(*service.HelloResp).Res)
		} else {
			fmt.Println(err)
		}
	}
	fmt.Println(send)
}
