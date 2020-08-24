
package loginClient

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"myRPC/client"
	"myRPC/meta"
	pbLogin "myRPC/myService/login/generate"
)

type loginClient struct {
}

var ClientCall = &loginClient{}

func NewClientCall() *loginClient {
	return &loginClient{}
}


func (c *loginClient)Login(ctx context.Context, req *pbLogin.LoginRequest, opts ...grpc.CallOption) (rsp *pbLogin.LoginResponse, err error) {
	clientObj,err := client.NewCommonClient()
	if err != nil {
		return
	}
	ctx = meta.InitClientMeta(ctx,clientObj.ServiceConf.ServiceName,"Login")
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
	defer conn.Close()
	newClient := pbLogin.NewLoginRpcClient(conn)
	return newClient.Login(ctx, req)
}

