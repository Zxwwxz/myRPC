
package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	pbHello "myRPC/myService/hello/generate"
	"myRPC/myService/hello/router"
	"net"
	"net/http"
)

func main()  {
	go promethusListen()
	grpcListen()
}

func promethusListen() {
	http.Handle("/metrics",promhttp.Handler())
	fmt.Println(http.ListenAndServe(":9091",nil))
}

func grpcListen() {
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
