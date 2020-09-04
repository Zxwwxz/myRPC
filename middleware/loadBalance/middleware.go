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
			clientMeta.RemainNodes = nil
			for _,v := range clientMeta.AllNodes {
				clientMeta.RemainNodes = append(clientMeta.RemainNodes,v)
			}
			failCount := 0
			for {
				//达到最大重连次数
				if failCount >= clientMeta.MaxReconnectNum{
					return nil,rpcConst.MaxReconnectFailed
				}
				//剩余节点中根据负载均衡找到一个节点
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
					//移除准备调用的节点
					clientMeta.RemainNodes = append(clientMeta.RemainNodes[:index], clientMeta.RemainNodes[index+1:]...)
				}
				logBase.Debug("LoadBalanceMiddleware,remainNodes=%v",clientMeta.RemainNodes)
				logBase.Debug("LoadBalanceMiddleware,curNode=%v",clientMeta.CurNode)
				//开始调用
				resp, err = next(ctx, req)
				//清理正在调用的节点
				clientMeta.CurNode = nil
				//有错误，需要换个节点重试
				if err != nil {
					if err == rpcConst.ConnFailed {
						failCount = failCount + 1
						continue
					}
				}
				//成功
				//如果只需要调用单一节点，可以返回了
				if clientMeta.CallerType != meta.Caller_type_all{
					return
				}
				//如果要调用所有节点，但已经没有可以调用的节点了，可以结束了
				if clientMeta.CallerType == meta.Caller_type_all && len(clientMeta.RemainNodes) == 0 {
					return nil,nil
				}
			}
		}
	}
}
