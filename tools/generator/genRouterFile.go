package generator

var routerTemplateFile = `
package router

import (
	"context"
	"{{.ImportPreFix}}/controller"
	{{.Package.Name}} "{{.ImportPreFix}}/generate"
)

type Router struct {}

{{range .Rpc}}
func (s *Router){{.Name}}(ctx context.Context, req *{{$.Package.Name}}.{{.RequestType}}) (rsp *{{$.Package.Name}}.{{.ReturnsType}}, err error) {
	controller := &controller.Controller{}
	rsp,err = controller.Check{{.Name}}(ctx,req)
	if err != nil {
		return
	}
	rsp,err = controller.Run{{.Name}}(ctx,req)
	if err != nil {
		return
	}
	return rsp,nil
}
{{end}}
`