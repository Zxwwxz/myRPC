
package router

import (
	"context"
	"myRPC/myService/hello/controller"
	pbHello "myRPC/myService/hello/generate"
)

type Router struct {}


func (s *Router)SayHello(ctx context.Context, req *pbHello.HelloRequest) (rsp *pbHello.HelloResponse, err error) {
	controller := &controller.Controller{}
	rsp,err = controller.CheckSayHello(ctx,req)
	if err != nil {
		return
	}
	rsp,err = controller.RunSayHello(ctx,req)
	if err != nil {
		return
	}
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

