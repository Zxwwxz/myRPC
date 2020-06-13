package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"myRPC/demo/grpc/pbHello"
)

func main()  {
	conn,err := grpc.Dial("127.0.0.1:8888",grpc.WithInsecure())
	if err != nil {
		fmt.Println("connect err:",err)
	}
	defer conn.Close()
	c := pbHello.NewHelloServiceClient(conn)
	resp,err := c.SayHello(context.TODO(),&pbHello.HelloRequest{Name:"zzz"})
	if err != nil {
		fmt.Println("rpc call err:",err)
	}
	fmt.Println("resp:",resp.Reply)
}
