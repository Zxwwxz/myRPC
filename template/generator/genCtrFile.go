package generator

var ctrTemplateFile = `
package controller

import "myRPC/service"

type Controller struct {
    CommonService *service.CommonService
}

func NewController(commonService *service.CommonService) *Controller {
    return &Controller{commonService}
}

func (controller *Controller)SetCommonService(commonService *service.CommonService) {
    controller.CommonService = commonService
}

func (controller *Controller)GetCommonService()(commonService *service.CommonService) {
    return controller.CommonService
}

func (controller *Controller)BeforeRpcLogic() {
}

func (controller *Controller)AfterRpcLogic() {
}
`

var ctrTemplateFuncFile = `
package controller

import (
	"context"
	{{.Package.Name}} "{{.ImportPreFix}}/proto"
)

{{range .Rpc}}
func (c *Controller)Run{{.Name}}(ctx context.Context, req *{{$.Package.Name}}.{{.RequestType}}) (rsp *{{$.Package.Name}}.{{.ReturnsType}}, err error) {
	rsp = new({{$.Package.Name}}.{{.ReturnsType}})
	return rsp,nil
}
{{end}}
`

var ctrTemplateStreamFuncFile = `
package controller

import (
	{{.Package.Name}} "{{.ImportPreFix}}/proto"
)

{{range .Rpc}}
func (c *Controller)Run{{.Name}}(req {{$.Package.Name}}.{{$.Service.Name}}_{{.Name}}Server) (err error) {
	return nil
}
{{end}}
`