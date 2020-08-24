
package controller

import (
	"context"
	pbLogin "myRPC/myService/login/generate"
)

type Controller struct {}


func (c *Controller)CheckLogin(ctx context.Context, req *pbLogin.LoginRequest) (rsp *pbLogin.LoginResponse, err error) {
	return rsp,nil
}

func (c *Controller)RunLogin(ctx context.Context, req *pbLogin.LoginRequest) (rsp *pbLogin.LoginResponse, err error) {
	rsp = new(pbLogin.LoginResponse)
	return rsp,nil
}

