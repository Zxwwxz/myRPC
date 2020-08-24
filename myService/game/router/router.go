
package router

import (
	"context"
	"myRPC/myService/game/controller"
	pbGame "myRPC/myService/game/generate"
	"myRPC/service"
)

type Router struct {}


func (s *Router)StartGame(ctx context.Context, req *pbGame.StartGameRequest) (rsp *pbGame.StartGameResponse, err error) {
	ctx,err = service.InitServiceFunc(ctx,"StartGame")
	if err != nil {
		return rsp,err
	}
	resultMwFunc := service.BuildServerMiddleware(s.MwFuncStartGame,nil,nil)
	newRsp,err := resultMwFunc(ctx,req)
	rsp = newRsp.(*pbGame.StartGameResponse)
	return rsp,err
}

func (s *Router)MwFuncStartGame(ctx context.Context, req interface{}) (rsp interface{}, err error) {
	newReq := req.(*pbGame.StartGameRequest)
	serverController := &controller.Controller{}
	//修改
	//rsp,err = serverController.CheckStartGame(ctx,newReq)
	//if err != nil {
	//	return
	//}
	rsp,err = serverController.RunStartGame(ctx,newReq)
	if err != nil {
		return
	}
	return rsp,nil
}

