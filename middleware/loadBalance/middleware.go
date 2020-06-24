package mwLoadBalance

import (
	"context"
	rpcConst "myRPC/const"
	balanceBase "myRPC/loadBalance"
	"myRPC/meta"
	mwBase "myRPC/middleware/base"
)

func LoadBalanceMiddleware(balancer balanceBase.BalanceInterface) mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			clientMeta := meta.GetClientMeta(ctx)
			if len(clientMeta.AllNodes) == 0 {
				err = rpcConst.NotHaveInstance
				return
			}
			for {
				clientMeta.CurNode, err = balancer.SelectNode(ctx, clientMeta.RemainNodes)
				if err != nil {
					return
				}
				delete(clientMeta.RemainNodes,clientMeta.CurNode.NodeId)
				resp, err = next(ctx, req)
				if err != nil {
					if err == rpcConst.ConnFailed {
						continue
					}
					return
				}
				break
			}
			return
		}
	}
}
