package wmConn

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"myRPC/client"
	"myRPC/meta"
	mwBase "myRPC/middleware/base"
)

func ConnMiddleware() mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			clientMeta := meta.GetClientMeta(ctx)
			if clientMeta.CurNode == nil{
				err = client.InvalidNode
				return
			}
			address := fmt.Sprintf("%s:%d", clientMeta.CurNode.NodeIp, clientMeta.CurNode.NodePort)
			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				return nil, client.ConnFailed
			}
			clientMeta.Conn = conn
			defer conn.Close()
			resp, err = next(ctx, req)
			return
		}
	}
}

