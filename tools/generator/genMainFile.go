package generator

var mainTemplateFile = `
package main

import (
	{{.Package.Name}} "{{.ImportPreFix}}/generate"
	"{{.ImportPreFix}}/router"
	"myRPC/service"
)

func main()  {
	service.Init()
	pbHello.RegisterHelloServiceServer(service.GetGrpcService(),&router.Router{})
	service.Run()
}

`
