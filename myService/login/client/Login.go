
package loginClient

import (
	"context"
	"errors"
	"myRPC/client"
	"myRPC/meta"
	pbLogin "myRPC/myService/login/generate"
)


func (c *loginClient)Login(ctx context.Context, req *pbLogin.LoginRequest, options []meta.ClientMetaOption) (rsp *pbLogin.LoginResponse, err error) {
	ctx,clientObj,err := client.InitClient(ctx,"login","Login",options)
	if err != nil {
		return
	}
	middlewareFunc := clientObj.BuildClientMiddleware(c.mwFuncLogin,nil,nil)
	newRsp, err := middlewareFunc(ctx, req)
	if err != nil {
		return nil, err
	}
	rsp = newRsp.(*pbLogin.LoginResponse)
	return rsp, err
}

func (c *loginClient)mwFuncLogin(ctx context.Context, request interface{}) (resp interface{}, err error) {
	clientMeta := meta.GetClientMeta(ctx)
	conn := clientMeta.Conn
	if conn == nil {
		return nil, errors.New("conn nil")
	}
	req := request.(*pbLogin.LoginRequest)
	newClient := pbLogin.NewLoginRpcClient(conn)
	return newClient.Login(ctx, req)
}

