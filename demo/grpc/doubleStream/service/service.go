package main

import (
	"fmt"
	"google.golang.org/grpc"
	"math/rand"
	"myRPC/demo/grpc/doubleStream/pbChat"
	"net"
	"time"
)

type server struct {}

func (s *server)Chat(req pbChat.ChatServer_ChatServer) error{
	rand.Seed(time.Now().Unix()+1)
	//不断接受客户端消息
	for {
		reqStruct,err := req.Recv()
		//客户关闭发送时，就会有err=EOF
		if err != nil {
			fmt.Println("server Recv err:", err)
			break
		}
		fmt.Println("server reqStruct:",reqStruct)
		rspStruct := &pbChat.RspStruct{
			RspMsg:fmt.Sprintf("%s,rsp msg:%d",reqStruct.ReqMsg,rand.Int()),
		}
		//给客户端回复消息
		err = req.Send(rspStruct)
		if err != nil {
			fmt.Println("server Send err:", err)
			break
		}
	}
	rspStruct := &pbChat.RspStruct{
		RspMsg:"server finish",
	}
	//即使客户端已经关闭发送，但客户端也可以接受消息
	err := req.Send(rspStruct)
	if err != nil {
		fmt.Println("server Send err:", err)
	}
	fmt.Println("server finish")
	time.Sleep(2*time.Second)
	//一旦执行完成，客户端接受消息处就会有err=EOF
	return nil
}

func main()  {
	//创建监听
	listen,err := net.Listen("tcp",":8889")
	if err != nil {
		fmt.Println("listen err:",err)
	}
	//创建grpc服务
	s := grpc.NewServer()
	//注册服务方法
	pbChat.RegisterChatServerServer(s,&server{})
	//开启服务
	err = s.Serve(listen)
	if err != nil {
		fmt.Println("start service err:",err)
	}
}
