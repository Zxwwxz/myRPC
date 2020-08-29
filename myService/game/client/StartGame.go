
package gameClient

import (
	"context"
	"errors"
	"myRPC/client"
	"myRPC/meta"
	pbGame "myRPC/myService/game/generate"
)


func (c *gameClient)StartGame(ctx context.Context, req *pbGame.StartGameRequest, options []meta.ClientMetaOption) (rsp *pbGame.StartGameResponse, err error) {
	ctx,clientObj,err := client.InitClient(ctx,"game","StartGame",options)
	if err != nil {
		return
	}
	middlewareFunc := clientObj.BuildClientMiddleware(c.mwFuncStartGame,nil,nil)
	newRsp, err := middlewareFunc(ctx, req)
	if err != nil {
		return nil, err
	}
	rsp = newRsp.(*pbGame.StartGameResponse)
	return rsp, err
}

func (c *gameClient)mwFuncStartGame(ctx context.Context, request interface{}) (resp interface{}, err error) {
	clientMeta := meta.GetClientMeta(ctx)
	conn := clientMeta.Conn
	if conn == nil {
		return nil, errors.New("conn nil")
	}
	req := request.(*pbGame.StartGameRequest)
	newClient := pbGame.NewGameServiceClient(conn)
	return newClient.StartGame(ctx, req)
}

