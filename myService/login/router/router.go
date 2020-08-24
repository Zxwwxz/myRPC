
package router

import (
	"context"
	"myRPC/meta"
	"myRPC/myService/login/controller"
	pbLogin "myRPC/myService/login/generate"
	"myRPC/service"
)

type Router struct {}


func (s *Router)Login(ctx context.Context, req *pbLogin.LoginRequest) (rsp *pbLogin.LoginResponse, err error) {
	serverMeta := &meta.ServerMeta{
		ServiceName:"pbLogin",
		ServiceMethod:"Login",
	}
	ctx = meta.SetServerMeta(ctx,serverMeta)
	resultMwFunc := service.BuildServerMiddleware(s.MwFuncLogin,nil,nil)
	newRsp,err := resultMwFunc(ctx,req)
	rsp = newRsp.(*pbLogin.LoginResponse)
	return rsp,err
}

func (s *Router)MwFuncLogin(ctx context.Context, req interface{}) (rsp interface{}, err error) {
	newReq := req.(*pbLogin.LoginRequest)
	serverController := &controller.Controller{}
	rsp,err = serverController.CheckLogin(ctx,newReq)
	if err != nil {
		return
	}
	rsp,err = serverController.RunLogin(ctx,newReq)
	if err != nil {
		return
	}
	return rsp,nil
}

