package generator

var routerTemplateFile = `
package router

import (
    "myRPC/myService/{{.ServiceName}}/controller"
)

type Router struct {
    Controller *controller.Controller
}

func NewRouter(controller *controller.Controller) *Router {
    return &Router{Controller:controller}
}

func (router *Router)SetController(controller *controller.Controller)()  {
    router.Controller = controller
}

func (router *Router)GetController()(*controller.Controller)  {
    return router.Controller
}
`

var routerTemplateFuncFile = `
package router

import (
	"context"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
	{{.Package.Name}} "{{.ImportPreFix}}/proto"
)

{{range .Rpc}}
func (s *Router){{.Name}}(ctx context.Context, req *{{$.Package.Name}}.{{.RequestType}}) (rsp *{{$.Package.Name}}.{{.ReturnsType}}, err error) {
	ctx,err = s.GetController().GetCommonService().InitServiceFunc(ctx,"{{.Name}}")
	if err != nil {
		return rsp,status.Error(codes.InvalidArgument, err.Error())
	}
	resultMwFunc := s.GetController().GetCommonService().BuildServerMiddleware(s.MwFunc{{.Name}},nil,nil)
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
	rsp,err = s.GetController().Run{{.Name}}(ctx,newReq)
	if err != nil {
		return
	}
	return rsp,nil
}
{{end}}
`

var routerTemplateStreamFuncFile = `
package router

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	{{.Package.Name}} "{{.ImportPreFix}}/proto"
)

{{range .Rpc}}
func (s *Router){{.Name}}(req {{$.Package.Name}}.{{$.Service.Name}}_{{.Name}}Server) (err error) {
	ctx,err := s.GetController().GetCommonService().InitServiceFunc(req.Context(),"{{.Name}}")
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	resultMwFunc := s.GetController().GetCommonService().BuildServerMiddleware(s.MwFunc{{.Name}},nil,nil)
	_,err = resultMwFunc(ctx,req)
    if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}

func (s *Router)MwFunc{{.Name}}(ctx context.Context, req interface{}) (rsp interface{}, err error) {
	newReq := req.({{$.Package.Name}}.{{$.Service.Name}}_{{.Name}}Server)
	err = s.GetController().Run{{.Name}}(newReq)
	if err != nil {
		return
	}
	return nil,nil
}
{{end}}
`