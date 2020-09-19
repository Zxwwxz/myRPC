package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"math/rand"
	"myRPC/demo/grpc/doubleStream/pbChat"
	"time"
)

func main()  {
	//创建连接
	conn, err := grpc.Dial("127.0.0.1:8889", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Dial err:", err)
		return
	}
	defer conn.Close()
	//得到客户端下的访问服务对象
	client := pbChat.NewChatServerClient(conn)
	//得到客户端下的访问服务方法的请求流对象
	chatObj,err := client.Chat(context.TODO())
	//开协程不断接受服务端消息
	go func() {
		for {
			rspStruct,err := chatObj.Recv()
			//即使发送已经关闭，也可以继续接受数据
			//服务端方法执行完成，就会有err=EOF
			if err != nil {
				fmt.Println("client Recv err:", err)
				return
			}
			fmt.Println("client resp:",rspStruct.RspMsg)
		}
	}()
	rand.Seed(time.Now().Unix())
	//发送消息
	for i:=1;i<=5;i++ {
		req := &pbChat.ReqStruct{
			ReqMsg:fmt.Sprintf("req msg:%d",rand.Int()),
		}
		err := chatObj.Send(req)
		if err != nil {
			fmt.Println("client Send err:", err)
			return
		}
	}
	//发送完成要关闭流对象，此时还是可以继续接收数据
	_ = chatObj.CloseSend()
	fmt.Println("client send finish")
	time.Sleep(5*time.Second)
}
