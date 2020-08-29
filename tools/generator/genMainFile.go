package generator

var mainTemplateFile = `
package main

import (
	"myRPC/log/base"
    "myRPC/service"
	"myRPC/myService/{{.ServiceName}}/generate"
	"myRPC/myService/{{.ServiceName}}/router"
	
)

func main()  {
	err := service.InitService()
	if err != nil {
		logBase.Fatal("init server err:%v",err)
	}
	{{.Package.Name}}.Register{{$.Service.Name}}Server(service.GetGrpcService(),&router.Router{})
	logBase.Debug("init server run")
	service.Run()
}

`
