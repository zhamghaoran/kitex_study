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

### 异常信息返回
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