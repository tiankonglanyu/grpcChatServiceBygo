# grpcChatServiceBygo
this is a chatservice by grpc and go demo, just for grpc learning
## 设置go module 开发模式（参考： https://github.com/golang/go/wiki/Modules）
## 生成客户端和服务端运行文件
- linux 请设置 go 的交叉编译环境
`SET CGO_ENABLED=0`
`SET GOOS=linux`
`SET GOARCH=amd64`
- 编译
```
cd grpcChatServiceBygo
go build -o rpc_server server.go   # 服户端
go build -o rpc_client client.go  # 客户端
```

![生成 客户端和服务端编译程序](https://user-images.githubusercontent.com/29748072/121644384-82fe3d00-cac5-11eb-8b94-9e360ade1c7c.png)


## 效果
- 服务端效果
![服务端效果](https://user-images.githubusercontent.com/29748072/121644545-b8a32600-cac5-11eb-9b45-dd7f639e4fe6.png)
- 启动两个客户端
![客户端 1](https://user-images.githubusercontent.com/29748072/121644703-ebe5b500-cac5-11eb-9dbd-23e293b02b6e.png)
![客户端 2](https://user-images.githubusercontent.com/29748072/121644721-f43df000-cac5-11eb-9758-ee0b75c521c5.png)
