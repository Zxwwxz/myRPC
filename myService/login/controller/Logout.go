
package controller

import (
	"context"
	pbLogin "myRPC/myService/login/generate"
)


func (c *Controller)RunLogout(ctx context.Context, req *pbLogin.LogoutRequest) (rsp *pbLogin.LogoutResponse, err error) {
	rsp = new(pbLogin.LogoutResponse)
	return rsp,nil
}

