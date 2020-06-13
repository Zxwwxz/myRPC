package generator

var routerTemplateFile = `
package router

import (
	"context"
	"myRPC/meta"
	"{{.ImportPreFix}}/controller"
	{{.Package.Name}} "{{.ImportPreFix}}/generate"
	"myRPC/service"
)

type Router struct {}

{{range .Rpc}}
func (s *Router){{.Name}}(ctx context.Context, req *{{$.Package.Name}}.{{.RequestType}}) (rsp *{{$.Package.Name}}.{{.ReturnsType}}, err error) {
	serverMeta := &meta.ServerMeta{
		ServiceName:"{{$.Package.Name}}",
		ServiceMethod:"{{.Name}}",
	}
	ctx = meta.SetServerMeta(ctx,serverMeta)
	resultMwFunc := service.BuildServerMiddleware(s.MwFunc{{.Name}},nil,nil)
	newRsp,err := resultMwFunc(ctx,req)
	rsp = newRsp.(*{{$.Package.Name}}.{{.ReturnsType}})
	return rsp,err
}

func (s *Router)MwFunc{{.Name}}(ctx context.Context, req interface{}) (rsp interface{}, err error) {
	newReq := req.(*{{$.Package.Name}}.{{.RequestType}})
	serverController := &controller.Controller{}
	rsp,err = serverController.Check{{.Name}}(ctx,newReq)
	if err != nil {
		return
	}
	rsp,err = serverController.Run{{.Name}}(ctx,newReq)
	if err != nil {
		return
	}
	return rsp,nil
}
{{end}}
`