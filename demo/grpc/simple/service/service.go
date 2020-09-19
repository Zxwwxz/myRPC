package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"myRPC/demo/grpc/simple/pbGuide"
	"net"
)

type server struct {}

func (s *server)GetFeature(ctx context.Context,req *pbGuide.Point) (rsp *pbGuide.Feature,err error) {
	return &pbGuide.Feature{Name:"simple resp",Location:req},nil
}

func main()  {
	listen,err := net.Listen("tcp",":8889")
	if err != nil {
		fmt.Println("listen err:",err)
	}
	s := grpc.NewServer()
	pbGuide.RegisterRouteGuideServer(s,&server{})
	err = s.Serve(listen)
	if err != nil {
		fmt.Println("start service err:",err)
	}
}
