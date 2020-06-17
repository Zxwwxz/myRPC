package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"myRPC/demo/grpc/pbHello"
	"net"
)

type server struct {}

func (s *server)SayHello(ctx context.Context, req *pbHello.HelloRequest) (rsp *pbHello.HelloResponse, err error) {
	md ,ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Println("md:",md)
	}
	rsp = &pbHello.HelloResponse{}
	rsp.Reply = req.Name+" good"
	return rsp,nil
}

func main()  {
	listen,err := net.Listen("tcp",":8889")
	if err != nil {
		fmt.Println("listen err:",err)
	}
	s := grpc.NewServer()
	pbHello.RegisterHelloServiceServer(s,&server{})
	err = s.Serve(listen)
	if err != nil {
		fmt.Println("start service err:",err)
	}
}
