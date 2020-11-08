package generator

var mainTemplateFile = `
package main

import (
	"myRPC/log/base"
	"myRPC/myService/{{.ServiceName}}/client"
	"myRPC/myService/{{.ServiceName}}/controller"
	"myRPC/myService/{{.ServiceName}}/proto"
	"myRPC/myService/{{.ServiceName}}/router"
	"myRPC/service"
)

func main()  {
	{{.ServiceName}}Controller := controller.NewController()
	commonService,err := service.NewService()
	if err != nil {
		logBase.Fatal("new server err:%v",err)
		return
	}
	{{.ServiceName}}Controller.SetCommonService(commonService)
	{{.ServiceName}}Client.NewClientCaller(commonService)
	{{.Package.Name}}.Register{{$.Service.Name}}Server(commonService.Server,&router.Router{ {{.ServiceName}}Controller })
	go func() {
		commonService.Stop()
	}()
	commonService.Run()
}
`
