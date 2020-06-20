package generator
var clientTemplateFile = `
package {{.ServiceName}}Client

import (
	"context"
	"google.golang.org/grpc"
	"myRPC/client"
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
	middlewareFunc := client.BuildClientMiddleware(c.mwFunc{{.Name}},nil,nil)
	newRsp, err := middlewareFunc(ctx, req)
	if err != nil {
		return nil, err
	}
	rsp = newRsp.(*{{$.Package.Name}}.{{.ReturnsType}})
	return rsp, err
}

func (c *{{$.ServiceName}}Client)mwFunc{{.Name}}(ctx context.Context, request interface{}) (resp interface{}, err error) {
	conn, err := grpc.Dial("127.0.0.1:8889", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	req := request.(*{{$.Package.Name}}.{{.RequestType}})
	defer conn.Close()
	newClient := {{$.Package.Name}}.New{{$.Service.Name}}Client(conn)
	return newClient.{{.Name}}(ctx, req)
}
{{end}}
`

