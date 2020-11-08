package generator
var clientTemplateFile = `
package {{.ServiceName}}Client

import (
	"myRPC/service"
)

var ClientCaller *{{.ServiceName}}Client

type {{.ServiceName}}Client struct {
	CommonService *service.CommonService
}

func NewClientCaller(commonService *service.CommonService) {
	if ClientCaller == nil {
		ClientCaller = &{{.ServiceName}}Client{commonService}
	}
}

func GetClientCaller() (*{{.ServiceName}}Client) {
	return ClientCaller
}

func (client *{{.ServiceName}}Client)SetCommonService(commonService *service.CommonService) {
	client.CommonService = commonService
}

func (client *{{.ServiceName}}Client)GetCommonService()(commonService *service.CommonService) {
	return client.CommonService
}
`

var clientTemplateFuncFile = `
package {{.ServiceName}}Client

import (
	"context"
	"errors"
	"myRPC/client"
    "myRPC/const"
	"myRPC/meta"
	"{{.ImportPreFix}}/proto"
)

{{range .Rpc}}
func (c *{{$.ServiceName}}Client){{.Name}}(ctx context.Context, req *{{$.Package.Name}}.{{.RequestType}}, options ...meta.ClientMetaOption) (rsp *{{$.Package.Name}}.{{.ReturnsType}}, err error) {
	ctx,clientObj,err := client.NewClient(ctx,"{{$.ServiceName}}","{{.Name}}",meta.Caller_mode_simple,c.GetCommonService().GetServiceConf(),options)
	if err != nil {
		newErr := rpcConst.ClientInitFailed
		newErr.Message = err.Error()
		return nil,newErr
	}
	middlewareFunc := clientObj.BuildClientMiddleware(c.mwFunc{{.Name}},nil,nil)
	newRsp, err := middlewareFunc(ctx, req)
	if err != nil {
		return nil, err
	}
	rsp,ok := newRsp.(*{{$.Package.Name}}.{{.ReturnsType}})
    if ok == false {
		return nil,nil
	}
	return rsp, nil
}

func (c *{{$.ServiceName}}Client)mwFunc{{.Name}}(ctx context.Context, req interface{}) (rsp interface{}, err error) {
	clientMeta := meta.GetClientMeta(ctx)
	conn := clientMeta.Conn
	if conn == nil {
		return nil, errors.New("conn nil")
	}
	newReq := req.(*{{$.Package.Name}}.{{.RequestType}})
	newClient := {{$.Package.Name}}.New{{$.Service.Name}}Client(conn)
	return newClient.{{.Name}}(ctx, newReq)
}
{{end}}
`

var clientTemplateStreamFuncFile = `
package {{.ServiceName}}Client

import (
	"context"
	"errors"
	"myRPC/client"
	"myRPC/const"
	"myRPC/meta"
	{{.Package.Name}} "{{.ImportPreFix}}/proto"
)

{{range .Rpc}}
func (c *{{$.ServiceName}}Client){{.Name}}(ctx context.Context, options ...meta.ClientMetaOption) (err error) {
	ctx,clientObj,err := client.NewClient(ctx,"{{$.ServiceName}}","{{.Name}}",meta.Caller_mode_stream,c.GetCommonService().GetServiceConf(),options)
	if err != nil {
		newErr := rpcConst.ClientInitFailed
		newErr.Message = err.Error()
		return newErr
	}
	middlewareFunc := clientObj.BuildClientMiddleware(c.mwFunc{{.Name}},nil,nil)
	_, err = middlewareFunc(ctx, nil)
	return err
}

func (c *{{$.ServiceName}}Client)mwFunc{{.Name}}(ctx context.Context, req interface{}) (rsp interface{}, err error) {
	clientMeta := meta.GetClientMeta(ctx)
	conn := clientMeta.Conn
	if conn == nil {
		return nil, errors.New("conn nil")
	}
	newClient := {{$.Package.Name}}.New{{$.Service.Name}}Client(conn)
	rsp,err =  newClient.{{.Name}}(ctx)
    if err != nil {
		return nil,err
	}
	callerModeFunc := clientMeta.CallerModeFunc
	if callerModeFunc != nil {
		callerModeFunc(rsp)
	}
	return nil, nil
}
{{end}}
`
