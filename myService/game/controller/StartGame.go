
package controller

import (
	"context"
	"fmt"
	"myRPC/config"
	"myRPC/log/base"
	"myRPC/myService/game/generate"
)


func (c *Controller)RunStartGame(ctx context.Context, req *pbGame.StartGameRequest) (rsp *pbGame.StartGameResponse, err error) {
	rsp = new(pbGame.StartGameResponse)
	logBase.Debug("req:%v",req)
	rsp.StartCode = "0"
	sid := config.GetConf().Base.ServiceId
	rsp.StartMsg = fmt.Sprintf("sid:%d,msg:%s",sid,"success:"+req.StartName)
	return rsp,nil
}

