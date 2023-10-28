package main

import (
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	service "kitex_study/kitex_gen/kitex_gen/service/hello"
	"log"
)

func main() {
	registry, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:12380"})
	if err != nil {
		panic(err)
	}
	svr := service.NewServer(new(HelloImpl), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "Hello"}), server.WithRegistry(registry))
	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
