package mwLoadBalance

import (
	"context"
	"myRPC/const"
	"myRPC/loadBalance/balancer"
	"myRPC/log/base"
	"myRPC/meta"
	"myRPC/middleware/base"
)

func LoadBalanceMiddleware(balancer balancer.BalanceInterface) mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			clientMeta := meta.GetClientMeta(ctx)
			if len(clientMeta.AllNodes) == 0 {
				return nil,rpcConst.NotFoundNode
			}
			for _,v := range clientMeta.AllNodes {
				clientMeta.RemainNodes = append(clientMeta.RemainNodes,v)
			}
			failCount := 0
			for {
				if failCount >= clientMeta.MaxReconnectNum{
					return nil,rpcConst.AllNodeFailed
				}
				clientMeta.CurNode, err = balancer.SelectNode(ctx, clientMeta.RemainNodes,clientMeta.ServiceName)
				if err != nil {
					return nil,rpcConst.AllNodeFailed
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
				logBase.Debug("LoadBalanceMiddleware,remainNodes=%v",clientMeta.RemainNodes)
				logBase.Debug("LoadBalanceMiddleware,curNode=%v",clientMeta.CurNode)
				resp, err = next(ctx, req)
				if err != nil {
					if err == rpcConst.ConnFailed {
						failCount = failCount + 1
						continue
					}
				}
				if clientMeta.CallerType != meta.Caller_type_all{
					return
				}
			}
		}
	}
}
