
package gameClient

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"myRPC/client"
	"myRPC/meta"
	pbGame "myRPC/myService/game/generate"
)

type gameClient struct {
}

func NewClientCall() *gameClient {
	return &gameClient{}
}


func (c *gameClient)StartGame(ctx context.Context, req *pbGame.StartGameRequest, opts ...grpc.CallOption) (rsp *pbGame.StartGameResponse, err error) {
	clientObj,err := client.InitClient()
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
	defer conn.Close()
	newClient := pbGame.NewGameServiceClient(conn)
	return newClient.StartGame(ctx, req)
}

