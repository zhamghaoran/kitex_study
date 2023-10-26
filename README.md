## kitex 教程

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
