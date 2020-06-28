package mwLoadBalance

import (
	"context"
	"fmt"
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
			for _,v := range clientMeta.AllNodes {
				clientMeta.RemainNodes = append(clientMeta.RemainNodes,v)
			}
			for {
				clientMeta.CurNode, err = balancer.SelectNode(ctx, clientMeta.RemainNodes)
				if err != nil {
					return
				}
				index := -1
				for i,v := range clientMeta.RemainNodes {
					if v.NodeId == clientMeta.CurNode.NodeId {
						index = i
					}
				}
				if index != -1 {
					clientMeta.RemainNodes = append(clientMeta.RemainNodes[:index], clientMeta.RemainNodes[index+1:]...)
				}
				fmt.Println("进入负载均衡中间件：",clientMeta.CurNode)
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
