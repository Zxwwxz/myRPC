package generator

var mainTemplateFile = `
package main

import (
	"fmt"
	"google.golang.org/grpc"
	{{.Package.Name}} "{{.ImportPreFix}}/generate"
	"{{.ImportPreFix}}/router"
	"net"
)

func main()  {
	listen,err := net.Listen("tcp",":8889")
	if err != nil {
		fmt.Println("listen err:",err)
	}
	s := grpc.NewServer()
	pbHello.Register{{.Service.Name}}Server(s,&router.Router{})
	err = s.Serve(listen)
	if err != nil {
		fmt.Println("start service err:",err)
	}
}
`
