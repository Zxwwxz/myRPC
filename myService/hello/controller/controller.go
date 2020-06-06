
package controller

import (
	"context"
	pbHello "myRPC/myService/hello/generate"
)

type Controller struct {}


func (c *Controller)CheckSayHello(ctx context.Context, req *pbHello.HelloRequest) (rsp *pbHello.HelloResponse, err error) {
	return rsp,nil
}

func (c *Controller)RunSayHello(ctx context.Context, req *pbHello.HelloRequest) (rsp *pbHello.HelloResponse, err error) {
	rsp = &pbHello.HelloResponse{}
	rsp.Reply = "reply:" + req.Name
	return rsp,nil
}

func (c *Controller)CheckSayBye(ctx context.Context, req *pbHello.ByeRequest) (rsp *pbHello.ByeResponse, err error) {
	return rsp,nil
}

func (c *Controller)RunSayBye(ctx context.Context, req *pbHello.ByeRequest) (rsp *pbHello.ByeResponse, err error) {
	return rsp,nil
}

