package generator

var testTemplateFile = `
package main

import (
    "context"
    "fmt"
    "myRPC/log/base"
    "myRPC/myService/{{.ServiceName}}/client"
    "myRPC/myService/{{.ServiceName}}/controller"
    "myRPC/myService/{{.ServiceName}}/proto"
    "myRPC/myService/{{.ServiceName}}/router"
    "myRPC/service"
    "time"
)

func InitRpc()(err error)  {
    commonService,err := service.NewService()
    if err != nil {
        logBase.Fatal("new server err:%v",err)
        return err
    }
    {{.ServiceName}}controller := controller.NewController(commonService)
    {{.ServiceName}}Client.NewClientCaller(commonService)
    {{.Package.Name}}.Register{{.Service.Name}}Server(commonService.Server,&router.Router{ {{.ServiceName}}controller} )
    go func() {
        commonService.Stop()
    }()
    go commonService.Run()
    time.Sleep(1*time.Second)
    return nil
}

{{range .Rpc}}
{{if eq .StreamsRequest true}}
func Test{{.Name}}() {
    err := {{$.ServiceName}}Client.GetClientCaller(nil).{{.Name}}(context.TODO(),meta.SetCallerModeFunc(func(i interface{}) {
        {{$.ServiceName}}Obj,ok := i.({{$.Package.Name}}.{{$.Service.Name}}_{{.Name}}Client)
        if ok == false {
            return
        }
        go func() {
            for {
                rspStruct,err := {{$.ServiceName}}Obj.Recv()
                if err != nil {
                    fmt.Println("Test{{.Name}} Recv err:", err)
                    return
                }
                fmt.Println("Test{{.Name}} Recv rsp:", rspStruct)
            }
        }()
        req := &{{$.Package.Name}}.{{.RequestType}}{}
        err := {{$.ServiceName}}Obj.Send(req)
        if err != nil {
            fmt.Println("Test{{.Name}},send,err:",err)
        }
        time.Sleep(2*time.Second)
    }))
    if err != nil {
        fmt.Println("Test{{.Name}},err:",err)
        return
    }
}
{{else if eq .StreamsRequest false}}
func Test{{.Name}}() {
    req := &{{$.Package.Name}}.{{.RequestType}}{}
    resp,err := {{$.ServiceName}}Client.GetClientCaller(nil).{{.Name}}(context.TODO(),req)
    if err != nil {
        fmt.Println("Test{{.Name}},err:",err)
        return
    }
    fmt.Println("Test{{.Name}},resp:",resp)
}
{{end}}
{{end}}

func main()  {
    err := InitRpc()
    if err != nil {
        fmt.Println("err:",err)
        return
    }
}
`
