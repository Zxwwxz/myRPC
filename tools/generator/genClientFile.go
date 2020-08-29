package generator
var clientTemplateFile = `
package {{.ServiceName}}Client

type {{.ServiceName}}Client struct {
}

func NewClientCall() *{{.ServiceName}}Client {
	return &{{.ServiceName}}Client{}
}
`

var clientTemplateFuncFile = `
package {{.ServiceName}}Client

import (
	"context"
	"errors"
	"myRPC/client"
	"myRPC/meta"
	{{.Package.Name}} "{{.ImportPreFix}}/generate"
)

{{range .Rpc}}
func (c *{{$.ServiceName}}Client){{.Name}}(ctx context.Context, req *{{$.Package.Name}}.{{.RequestType}}, options []meta.ClientMetaOption) (rsp *{{$.Package.Name}}.{{.ReturnsType}}, err error) {
	ctx,clientObj,err := client.InitClient(ctx,"{{$.ServiceName}}","{{.Name}}",options)
	if err != nil {
		return
	}
	middlewareFunc := clientObj.BuildClientMiddleware(c.mwFunc{{.Name}},nil,nil)
	newRsp, err := middlewareFunc(ctx, req)
	if err != nil {
		return nil, err
	}
	rsp = newRsp.(*{{$.Package.Name}}.{{.ReturnsType}})
	return rsp, err
}

func (c *{{$.ServiceName}}Client)mwFunc{{.Name}}(ctx context.Context, request interface{}) (resp interface{}, err error) {
	clientMeta := meta.GetClientMeta(ctx)
	conn := clientMeta.Conn
	if conn == nil {
		return nil, errors.New("conn nil")
	}
	req := request.(*{{$.Package.Name}}.{{.RequestType}})
	newClient := {{$.Package.Name}}.New{{$.Service.Name}}Client(conn)
	return newClient.{{.Name}}(ctx, req)
}
{{end}}
`
