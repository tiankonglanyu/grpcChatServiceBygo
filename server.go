// @Time: 2021/5/5 20:23
// @Author: yxl
// @File: server.go
package main

import (
	"google.golang.org/grpc"
	"grpc_chat/chatserver"
	"log"
	"net"
	"os"
)

func main() {
	// set the port

	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "5000" // 没有的话设置为5000
	}

	listen, err := net.Listen("tcp", ":" + Port)  // 开通一个 tcp 的监听
	if err != nil {
		log.Fatalf("count not listen port: %v, err is %v", Port, err)  // log.Fatal 会直接终止后面的程序
	}

	log.Printf("success listen port : %v", Port)  // log 打印会自动带有时间信息，比较方便

	grpcserver := grpc.NewServer()  // 初始化一个 gprc 服务对象

	cs := chatserver.ChatServer{} // 类似普通的 rpc, 实例化结构体，穿进去，他就可以使用这个命名空间的函数了
	chatserver.RegisterServicesServer(grpcserver, &cs)  // 结构体注册到 grpc.NewServer 实例
	err = grpcserver.Serve(listen)   // 服务正在监听的端口

	if err != nil {
		log.Fatalf("fail to start grpc server because %v", err)

	}

}
