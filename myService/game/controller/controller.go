
package controller

import (
	"context"
	pbGame "myRPC/myService/game/generate"
)

type Controller struct {}


func (c *Controller)CheckStartGame(ctx context.Context, req *pbGame.StartGameRequest) (rsp *pbGame.StartGameResponse, err error) {
	return rsp,nil
}

func (c *Controller)RunStartGame(ctx context.Context, req *pbGame.StartGameRequest) (rsp *pbGame.StartGameResponse, err error) {
	rsp = new(pbGame.StartGameResponse)
	return rsp,nil
}

