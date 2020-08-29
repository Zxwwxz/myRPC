
package loginClient

import (
	"context"
	"errors"
	"myRPC/client"
	"myRPC/meta"
	pbLogin "myRPC/myService/login/generate"
)


func (c *loginClient)Logout(ctx context.Context, req *pbLogin.LogoutRequest, options []meta.ClientMetaOption) (rsp *pbLogin.LogoutResponse, err error) {
	ctx,clientObj,err := client.InitClient(ctx,"login","Logout",options)
	if err != nil {
		return
	}
	middlewareFunc := clientObj.BuildClientMiddleware(c.mwFuncLogout,nil,nil)
	newRsp, err := middlewareFunc(ctx, req)
	if err != nil {
		return nil, err
	}
	rsp = newRsp.(*pbLogin.LogoutResponse)
	return rsp, err
}

func (c *loginClient)mwFuncLogout(ctx context.Context, request interface{}) (resp interface{}, err error) {
	clientMeta := meta.GetClientMeta(ctx)
	conn := clientMeta.Conn
	if conn == nil {
		return nil, errors.New("conn nil")
	}
	req := request.(*pbLogin.LogoutRequest)
	newClient := pbLogin.NewLoginRpcClient(conn)
	return newClient.Logout(ctx, req)
}

