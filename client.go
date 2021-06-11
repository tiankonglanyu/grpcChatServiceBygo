// @Time   : 2021/5/5 21:27
// @Author : yxl
// @File   : client.go

package main

import (
	"bufio"
	"context"
	"fmt"
	"grpc_chat/chatserver"
	"log"
	"os"
	"strings"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Enter Server IP:Port ::: ")
	reader := bufio.NewReader(os.Stdin)  // 控制台读取数据
	serverID, err := reader.ReadString('\n')  // 读取服务 地址

	if err != nil {
		log.Printf("Failed to read from console :: %v", err)
	}
	serverID = strings.Trim(serverID, "\r\n")

	log.Println("Connecting : " + serverID)

	//connect to grpc server
	conn, err := grpc.Dial(serverID, grpc.WithInsecure())  // 拨号请求 grpc 服务

	if err != nil {
		log.Fatalf("Faile to conncet to gRPC server :: %v", err)
	}
	defer conn.Close( )

	//call ChatService to create a stream
	client := chatserver.NewServicesClient(conn)  // 穿件 grpc 客户端

	stream, err := client.ChatService(context.Background())  // 调用客户端方法  返回 stream 流对象， 参见 客户端 生成代签名
	if err != nil {
		log.Fatalf("Failed to call ChatService :: %v", err)
	}

	// implement communication with gRPC server
	ch := clienthandle{stream: stream}
	ch.clientConfig()
	go ch.sendMessage()  // 发送消息
	go ch.receiveMessage()  // 接受消息

	//blocker
	bl := make(chan bool)
	<- bl

}

//clienthandle
type clienthandle struct {  // 消息体封装
	stream     chatserver.Services_ChatServiceClient
	clientName string
}

func (ch *clienthandle) clientConfig() {  // 写入 clientName

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Your Name : ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf(" Failed to read from console :: %v", err)
	}
	ch.clientName = strings.Trim(name, "\r\n")

}

//send message
func (ch *clienthandle) sendMessage() {

	// create a loop
	for {

		reader := bufio.NewReader(os.Stdin)
		clientMessage, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf(" Failed to read from console :: %v", err)
		}
		clientMessage = strings.Trim(clientMessage, "\r\n")

		clientMessageBox := &chatserver.FormClient{
			Name: ch.clientName,
			Body: clientMessage,
		}

		err = ch.stream.Send(clientMessageBox)

		if err != nil {
			log.Printf("Error while sending message to server :: %v", err)
		}

	}

}

//receive message
func (ch *clienthandle) receiveMessage() {

	//create a loop
	for {
		mssg, err := ch.stream.Recv()
		if err != nil {
			log.Printf("Error in receiving message from server :: %v", err)
		}

		//print message to console
		fmt.Printf("%s : %s \n", mssg.Name, mssg.Body)
	}
}
