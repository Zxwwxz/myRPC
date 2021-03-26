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

//连接中间件
func ConnMiddleware() mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			clientMeta := meta.GetClientMeta(ctx)
			if clientMeta.CurNode == nil{
				return nil,rpcConst.NotFoundNode
			}
			address := fmt.Sprintf("%s:%s", clientMeta.CurNode.NodeIp, clientMeta.CurNode.NodePort)
			logBase.Debug("ConnMiddleware,address=%s",address)
			//创建连接
			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				return nil, rpcConst.ConnFailed
			}
			//保存连接
			clientMeta.Conn = conn
			defer conn.Close()
			return next(ctx, req)
		}
	}
}

