package generator
var clientTemplateFile = `
package {{.ServiceName}}Client

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"myRPC/client"
	"myRPC/meta"
	{{.Package.Name}} "{{.ImportPreFix}}/generate"
)

type {{.ServiceName}}Client struct {
}

var ClientCall = &{{.ServiceName}}Client{}

func NewClientCall() *{{.ServiceName}}Client {
	return &{{.ServiceName}}Client{}
}

{{range .Rpc}}
func (c *{{$.ServiceName}}Client){{.Name}}(ctx context.Context, req *{{$.Package.Name}}.{{.RequestType}}, opts ...grpc.CallOption) (rsp *{{$.Package.Name}}.{{.ReturnsType}}, err error) {
	clientObj,err := client.NewCommonClient()
	if err != nil {
		return
	}
	middlewareFunc := clientObj.BuildClientMiddleware(c.mwFuncLogin,nil,nil)
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
	if conn != nil {
		return nil, errors.New("conn nil")
	}
	req := request.(*{{$.Package.Name}}.{{.RequestType}})
	defer conn.Close()
	newClient := {{$.Package.Name}}.New{{$.Service.Name}}Client(conn)
	return newClient.{{.Name}}(ctx, req)
}
{{end}}
`

