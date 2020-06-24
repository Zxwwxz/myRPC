package generator

var mainTemplateFile = `
package main

import (
	"fmt"
	{{.Package.Name}} "{{.ImportPreFix}}/generate"
	"{{.ImportPreFix}}/router"
	"myRPC/service"
)

func main()  {
	err := service.Init()
	if err != nil {
		fmt.Println("init server err:",err)
	}
	{{.Package.Name}}.Register{{$.Service.Name}}Server(service.GetGrpcService(),&router.Router{})
	fmt.Println("init server success")	
	service.Run()
}

`
