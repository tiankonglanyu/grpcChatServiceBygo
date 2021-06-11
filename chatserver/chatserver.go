// @Time   : 2021/5/5 20:38
// @Author : yxl
// @File   : chatserver.go
package chatserver

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type messageUint struct {  // 消息结构体
	ClientName        string
	MessageBody       string
	MessageUniqueCode int
	ClientUniqueCode  int
}

type messageHandle struct {   // 全局变量， 包含消息结构体切片和 锁
	MQue []messageUint
	mu   sync.Mutex
}

var messageHandleObj = messageHandle{}  // 初始化全局消息 和 锁 对象

type ChatServer struct {  // 实现服务的结构体 ，他里面需要实现全部的 Services  定义的 rpc 函数, 按照生成的 .pb.go 文件接口定义的函数签名来写
}

// define services
func (ch *ChatServer) ChatService(cc Services_ChatServiceServer) error {  // <Services name>rpc_nameServer
	client_code := rand.Intn(1e6)   // 生成一个 client code， 每个客户端调用 函数的时候，都会生成一个随机的  id, 填充在消息里面，标记
	errch := make(chan error)  // 创建一个 error 的管道

	// receive message
	go receive(cc, client_code, errch)

	// send message
	go send(cc, client_code, errch)

	return <- errch

}

// receive message
func receive(cc Services_ChatServiceServer, client_code int, errch chan error) {
	// loop for receive message
	for {
		mes, err := cc.Recv()
		if err != nil {
			errch <- err
			log.Printf("error when recv message, %v", err)

		} else {
			messageHandleObj.mu.Lock()
			messageHandleObj.MQue = append(messageHandleObj.MQue, messageUint{
				ClientName:        mes.Name,
				MessageBody:       mes.Body,
				MessageUniqueCode: rand.Intn(1e8),
				ClientUniqueCode:  client_code,
			})
			messageHandleObj.mu.Unlock()
			log.Printf("message num is %v", messageHandleObj.MQue[len(messageHandleObj.MQue)-1]) // 打印出最新的一条消息
		}
	}
}

// send message
func send(cc Services_ChatServiceServer, client_code int, errch chan error) {

	// loop messageHandleObj.MQue
	for { // 没消息也要继续监听管道的动态

		for {
			time.Sleep(time.Second)
			messageHandleObj.mu.Lock()
			if len(messageHandleObj.MQue) < 1 { // 如果没消息了，就终止发消息
				messageHandleObj.mu.Unlock()
				break
			}

			senderUniqueCode := messageHandleObj.MQue[0].ClientUniqueCode
			senderName4Client := messageHandleObj.MQue[0].ClientName
			message4Client := messageHandleObj.MQue[0].MessageBody

			messageHandleObj.mu.Unlock()

			if client_code != senderUniqueCode { // 不给自己发消息
				err := cc.Send(&FromServer{
					Name: senderName4Client,
					Body: message4Client,
				})
				if err != nil {
					errch <- err
				}
				messageHandleObj.mu.Lock()
				if len(messageHandleObj.MQue) > 1 {
					messageHandleObj.MQue = messageHandleObj.MQue[1:] // 发送过的消息删掉
				} else {
					messageHandleObj.MQue = []messageUint{}
				}
				messageHandleObj.mu.Unlock()

			}
		}
		time.Sleep(time.Second)

	}

}


/* 结构主体:
server 全局 定义消息容器， go 单独运行进城 服务每个客户端
一旦发现消息数据更新， 并且 不是当前客户端的 id ,就吧消息转发给他， ok, 客户端收到消息就 输出到控制台
prefect  练习

注意点： 因为是全局消息对象，所以操作上下需要加 锁， 如果是管道则不需要
*/
