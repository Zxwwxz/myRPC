
package router

import (
	"context"
	"myRPC/myService/login/controller"
	pbLogin "myRPC/myService/login/generate"
	"myRPC/service"
)


func (s *Router)Logout(ctx context.Context, req *pbLogin.LogoutRequest) (rsp *pbLogin.LogoutResponse, err error) {
	ctx,err = service.InitServiceFunc(ctx,"Logout")
	if err != nil {
		return rsp,err
	}
	resultMwFunc := service.BuildServerMiddleware(s.MwFuncLogout,nil,nil)
	newRsp,err := resultMwFunc(ctx,req)
	rsp = newRsp.(*pbLogin.LogoutResponse)
	return rsp,err
}

func (s *Router)MwFuncLogout(ctx context.Context, req interface{}) (rsp interface{}, err error) {
	newReq := req.(*pbLogin.LogoutRequest)
	serverController := &controller.Controller{}
	rsp,err = serverController.RunLogout(ctx,newReq)
	if err != nil {
		return
	}
	return rsp,nil
}

