package mwLoadBalance

import (
	"context"
	"github.com/ibinarytree/koala/errno"
	"github.com/ibinarytree/koala/loadbalance"
	"github.com/ibinarytree/koala/logs"
	"myRPC/client"
	balanceBase "myRPC/loadBalance/base"
	"myRPC/meta"
	mwBase "myRPC/middleware/base"
)

func LoadBalanceMiddleware(balancer balanceBase.BalanceInterface) mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			clientMeta := meta.GetClientMeta(ctx)
			if len(clientMeta.AllNodes) == 0 {
				err = client.NotHaveInstance
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
					if err == client.ConnFailed {
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
