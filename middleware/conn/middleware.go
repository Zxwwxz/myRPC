package wmConn

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	rpcConst "myRPC/const"
	"myRPC/log/base"
	"myRPC/meta"
	mwBase "myRPC/middleware/base"
)

func ConnMiddleware() mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			clientMeta := meta.GetClientMeta(ctx)
			if clientMeta.CurNode == nil{
				err = rpcConst.NotFoundNode
				return
			}
			address := fmt.Sprintf("%s:%s", clientMeta.CurNode.NodeIp, clientMeta.CurNode.NodePort)
			logBase.Debug("ConnMiddleware,address=%s",address)
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

