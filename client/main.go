package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/kitex/client"
	"kitex_study/kitex_gen/kitex_gen/service"
	"kitex_study/kitex_gen/kitex_gen/service/hello"
)

func main() {
	newClient, err := hello.NewClient("hello", client.WithHostPorts("0.0.0.0:8888"))
	if err != nil {
		return
	}
	req := &service.HelloReq{Name: "zhr"}
	send, err := newClient.Send(context.Background(), req)
	if err != nil {
		return
	}
	fmt.Println(send)
}
