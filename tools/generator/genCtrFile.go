package generator

var ctrTemplateFile = `
package controller

import (
	"context"
	{{.Package.Name}} "{{.ImportPreFix}}/generate"
)

type Controller struct {}

{{range .Rpc}}
func (c *Controller)Check{{.Name}}(ctx context.Context, req *{{$.Package.Name}}.{{.RequestType}}) (rsp *{{$.Package.Name}}.{{.ReturnsType}}, err error) {
	return rsp,nil
}

func (c *Controller)Run{{.Name}}(ctx context.Context, req *{{$.Package.Name}}.{{.RequestType}}) (rsp *{{$.Package.Name}}.{{.ReturnsType}}, err error) {
	return rsp,nil
}
{{end}}
`
