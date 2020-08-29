
package controller

import (
	"context"
	"myRPC/myService/login/generate"
)


func (c *Controller)RunLogin(ctx context.Context, req *pbLogin.LoginRequest) (rsp *pbLogin.LoginResponse, err error) {
	rsp = new(pbLogin.LoginResponse)
	return rsp,nil
}

