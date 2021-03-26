package generator

var mainTemplateFile = `
package main

import (
	"myRPC/log/base"
	"myRPC/myService/{{.ServiceName}}/controller"
	"myRPC/myService/{{.ServiceName}}/proto"
	"myRPC/myService/{{.ServiceName}}/router"
	"myRPC/service"
)

//服务入口
//服务流转：
//请求：A->B，A中直接导入B的client包的对应方法
//接受：一个请求过来，router->controller
func main()  {
	commonService,err := service.NewService()
	if err != nil {
		logBase.Fatal("new server err:%v",err)
		return
	}
    {{.ServiceName}}Controller := controller.NewController(commonService)
	{{.Package.Name}}.Register{{$.Service.Name}}Server(commonService.Server,router.NewRouter( {{.ServiceName}}Controller ))
	go func() {
		commonService.Stop()
	}()
	commonService.Run()
}
`
