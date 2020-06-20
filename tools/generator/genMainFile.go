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
	{{.Package.Name}}.Register{{$.Service.Name}}Server(service.GetGrpcService(),&router.Router{})
	service.Run()
}

`
