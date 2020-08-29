package generator

var ctrTemplateFile = `
package controller

type Controller struct {}
`

var ctrTemplateFuncFile = `
package controller

import (
	"context"
	{{.Package.Name}} "{{.ImportPreFix}}/generate"
)

{{range .Rpc}}
func (c *Controller)Run{{.Name}}(ctx context.Context, req *{{$.Package.Name}}.{{.RequestType}}) (rsp *{{$.Package.Name}}.{{.ReturnsType}}, err error) {
	rsp = new({{$.Package.Name}}.{{.ReturnsType}})
	return rsp,nil
}
{{end}}
`
