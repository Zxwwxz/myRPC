
package router

import (
	"context"
	"fmt"
	mwBase "myRPC/middleware/base"
	"myRPC/myService/hello/controller"
	pbHello "myRPC/myService/hello/generate"
)

type Router struct {}

func (s *Router)MwSayHello1(wareFunc mwBase.MiddleWareFunc) (mwBase.MiddleWareFunc) {
	return func(ctx context.Context, req interface{}) (rsp interface{}, err error) {
		fmt.Println("进入中间件1-》前")
		rsp,err = wareFunc(ctx,req)
		fmt.Println("进入中间件1-》后")
		return
	}
}

func (s *Router)MwSayHello2(wareFunc mwBase.MiddleWareFunc) (mwBase.MiddleWareFunc) {
	return func(ctx context.Context, req interface{}) (rsp interface{}, err error) {
		fmt.Println("进入中间件2-》前")
		rsp,err = wareFunc(ctx,req)
		fmt.Println("进入中间件2-》后")
		return
	}
}

func (s *Router)SayHello(ctx context.Context, req *pbHello.HelloRequest) (rsp *pbHello.HelloResponse, err error) {
	resultMw := mwBase.Chain(s.MwSayHello1,s.MwSayHello2)
	resultMwFunc := resultMw(s.MwFuncSayHello)
	newRsp,err := resultMwFunc(ctx,req)
	rsp = newRsp.(*pbHello.HelloResponse)
	return rsp,err
}

func (s *Router)MwFuncSayHello(ctx context.Context, req interface{}) (rsp interface{}, err error) {
	fmt.Println("进入函数-》前")
	newReq := req.(*pbHello.HelloRequest)
	serverController := &controller.Controller{}
	rsp,err = serverController.CheckSayHello(ctx,newReq)
	if err != nil {
		return
	}
	rsp,err = serverController.RunSayHello(ctx,newReq)
	if err != nil {
		return
	}
	fmt.Println("进入函数-》后")
	return rsp,nil
}

func (s *Router)SayBye(ctx context.Context, req *pbHello.ByeRequest) (rsp *pbHello.ByeResponse, err error) {
	controller := &controller.Controller{}
	rsp,err = controller.CheckSayBye(ctx,req)
	if err != nil {
		return
	}
	rsp,err = controller.RunSayBye(ctx,req)
	if err != nil {
		return
	}
	return rsp,nil
}

