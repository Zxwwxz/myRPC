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
			address := fmt.Sprintf("%s:%s", "127.0.0.1", clientMeta.CurNode.NodePort)
			//创建连接
			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				return nil, rpcConst.ConnFailed
			}
			clientMeta.Conn = conn
			fmt.Println("进入连接中间件：",address)
			defer conn.Close()
			resp, err = next(ctx, req)
			return
		}
	}
}

