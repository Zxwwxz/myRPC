package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"myRPC/demo/grpc/pbHello"
)

func main()  {
	conn,err := grpc.Dial("127.0.0.1:8888",grpc.WithInsecure())
	if err != nil {
		fmt.Println("connect err:",err)
	}
	defer conn.Close()
	c := pbHello.NewHelloServiceClient(conn)
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx,"req_key","req_value")
	resp,err := c.SayHello(context.TODO(),&pbHello.HelloRequest{Name:"zzz"})
	if err != nil {
		fmt.Println("rpc call err:",err)
	}
	fmt.Println("resp:",resp.Reply)
}

//func main()  {
//	resp ,err := helloClient.ClientCall.SayHello(context.Background(),&pbHello.HelloRequest{Name:"aaa"})
//	if err != nil {
//		fmt.Println("rpc call err:",err)
//	}
//	fmt.Println("resp:",resp.Reply)
//}
