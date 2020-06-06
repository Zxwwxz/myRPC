
package main

import (
	"fmt"
	"google.golang.org/grpc"
	pbHello "myRPC/myService/hello/generate"
	"myRPC/myService/hello/router"
	"net"
)

func main()  {
	listen,err := net.Listen("tcp",":8889")
	if err != nil {
		fmt.Println("listen err:",err)
	}
	s := grpc.NewServer()
	pbHello.RegisterHelloServiceServer(s,&router.Router{})
	err = s.Serve(listen)
	if err != nil {
		fmt.Println("start service err:",err)
	}
}
