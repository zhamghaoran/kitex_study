## kitex 教程
- 官方文档：https://www.cloudwego.io/zh/docs/kitex/
## 简单demo

### 编写proto文件
```protobuf
syntax = "proto3";

option go_package = "kitex_gen/service";

service hello {
    rpc Send (HelloReq) returns (HelloResp);
}
message HelloReq {
    string name = 1;
}
message HelloResp {
    string res = 1;
}
```
#### 注意
    - proto 文件当中的go_package 应该以kitex_gen 开头
安装kitex 
```shell
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest

```
安装完成之后在终端输入kitex 查看是否有相应的输出
### 执行命令
```shell
kitex -module kitex_study -service hello-service service.proto
```
`-module` 中的名字应该与go.mod 中的ModuleName相同,`-service` 后跟上service的名字，最后跟上proto文件名。
执行完命令之后会自动创建好 main.go ,  handler.go 和 kitex_gen 文件夹。其中handler.go 当中为需要自己实现的接口逻辑

那么实现完接口逻辑之后服务端的程序就算完成了。
### client
创建client 文件夹
```go
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
```
完成客户端编写

## 异常信息返回
如果调用链路跑不通的话，肯定可以获取到err的值，但是如果我们想要获取到逻辑错误返回的err，以便于链路err做区分，那么就需要使用额外的方法来包装错误

服务端错误封装
```go
func (s *HelloImpl) Send(ctx context.Context, req *service.HelloReq) (resp *service.HelloResp, err error) {
	resp = &service.HelloResp{Res: "hello " + req.Name}
	err = kerrors.NewGRPCBizStatusError(500, "err")
	grpcStatusErr := err.(kerrors.GRPCStatusIface)
	st, _ := grpcStatusErr.GRPCStatus().WithDetails(resp)
	grpcStatusErr.SetGRPCStatus(st)
	return
}
```
客户端获取错误
```go
send, err := newClient.Send(context.Background(), req)
	if err != nil {
		if bizErr, ok := kerrors.FromBizStatusError(err); ok {
			fmt.Println(bizErr.BizStatusCode())
			fmt.Println(bizErr.BizMessage())
			// 通过类型断言获取到detail信息
			println(bizErr.(status.Iface).GRPCStatus().Details()[0].(*service.HelloResp).Res)
		}
	}
```
## 整合etcd 完成服务发现
### server 端代码
```go
package main

import (
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	service "kitex_study/kitex_gen/kitex_gen/service/hello"
	"log"
)

func main() {
	// 填写etcd 的地址
	registry, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:12380"})
	if err != nil {
		panic(err)
	}
	// 现在本地注册服务，然后再向etcd 注册服务
	svr := service.NewServer(new(HelloImpl), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "Hello"}), server.WithRegistry(registry))
	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
```
### client 端代码
同理，在client 当中指定etcd的地址，获取到etcdResolver，然后在获取连接客户端的时候指定服务的名称即可完成服务发现
```go
registry, err := etcd.NewEtcdResolver([]string{"127.0.0.1:12380"})
	if err != nil {
		panic(err)
	}
	newClient, err := hello.NewClient("Hello", client.WithResolver(registry), client.WithTransportProtocol(transport.GRPC))
```
## 超时控制
在方法调用的时候添加 `callopt.WithRPCTimeout(2 * time.Second)` 或者 在客户端配置当中加上`client.WithRPCTimeout(2 * time.Second)`
客户端中配置：
```go
newClient, err := hello.NewClient("Hello", client.WithResolver(registry), client.WithRPCTimeout(time.Second*2))
```
调用方法时配置：
```go
send, err := newClient.Send(context.Background(), req, callopt.WithRPCTimeout(time.Second*2))
```
可以根据具体场景来选择超时控制的力度

可以通过
```go
kerrors.IsTimeoutError(err)
```
来判断这个err是不是因为超时导致的
