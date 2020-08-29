package main

import (
	"context"
	"fmt"
	"myRPC/meta"
	"myRPC/myService/login/client"
	"myRPC/myService/login/generate"
)

//func main()  {
//	conn,err := grpc.Dial("127.0.0.1:8888",grpc.WithInsecure())
//	if err != nil {
//		fmt.Println("connect err:",err)
//	}
//	defer conn.Close()
//	c := pbHello.NewHelloServiceClient(conn)
//	ctx := context.Background()
//	ctx = metadata.AppendToOutgoingContext(ctx,"req_key","req_value")
//	resp,err := c.SayHello(context.TODO(),&pbHello.HelloRequest{Name:"zzz"})
//	if err != nil {
//		fmt.Println("rpc call err:",err)
//	}
//	fmt.Println("resp:",resp.Reply)
//}

func main()  {
	resp ,err := loginClient.NewClientCall().Login(context.Background(),&pbLogin.LoginRequest{LoginName:"aaa"},[]meta.ClientMetaOption{
		meta.SetMaxReconnectNum(3),
		meta.SetCallerType(meta.Caller_type_balance),
	})
	if err != nil {
		fmt.Println("rpc call err:",err)
	}
	fmt.Println("resp:",resp)
}
