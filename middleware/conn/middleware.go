package wmConn

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	rpcConst "myRPC/const"
	"myRPC/meta"
	mwBase "myRPC/middleware/base"
)

func ConnMiddleware() mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			clientMeta := meta.GetClientMeta(ctx)
			if clientMeta.CurNode == nil{
				err = rpcConst.InvalidNode
				return
			}
			address := fmt.Sprintf("%s:%d", clientMeta.CurNode.NodeIp, clientMeta.CurNode.NodePort)
			//创建连接
			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				return nil, rpcConst.ConnFailed
			}
			clientMeta.Conn = conn
			defer conn.Close()
			resp, err = next(ctx, req)
			return
		}
	}
}

