package generator

var routerTemplateFile = `
package router

type Router struct {}
`

var routerTemplateFuncFile = `
package router

import (
	"context"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
	"{{.ImportPreFix}}/controller"
	{{.Package.Name}} "{{.ImportPreFix}}/generate"
	"myRPC/service"
)

{{range .Rpc}}
func (s *Router){{.Name}}(ctx context.Context, req *{{$.Package.Name}}.{{.RequestType}}) (rsp *{{$.Package.Name}}.{{.ReturnsType}}, err error) {
	ctx,err = service.InitServiceFunc(ctx,"{{.Name}}")
	if err != nil {
		return rsp,status.Error(codes.InvalidArgument, err.Error())
	}
	resultMwFunc := service.BuildServerMiddleware(s.MwFunc{{.Name}},nil,nil)
	newRsp,err := resultMwFunc(ctx,req)
    if err != nil {
		return nil,status.Error(codes.Internal, err.Error())
	}
	rsp,ok := newRsp.(*{{$.Package.Name}}.{{.ReturnsType}})
	if ok == false {
		return nil,status.Error(codes.InvalidArgument, "rsp type illegal")
	}
	return rsp,nil
}

func (s *Router)MwFunc{{.Name}}(ctx context.Context, req interface{}) (rsp interface{}, err error) {
	newReq := req.(*{{$.Package.Name}}.{{.RequestType}})
	serverController := &controller.Controller{}
	rsp,err = serverController.Run{{.Name}}(ctx,newReq)
	if err != nil {
		return
	}
	return rsp,nil
}
{{end}}
`