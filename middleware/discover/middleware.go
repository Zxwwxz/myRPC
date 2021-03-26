package mwDiscover

import (
	"context"
	"myRPC/const"
	"myRPC/log/base"
	"myRPC/meta"
	"myRPC/middleware/base"
	"myRPC/registry/register"
)

//服务发现中间件
func DiscoveryMiddleware(discovery register.RegisterInterface) mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			clientMeta := meta.GetClientMeta(ctx)
			//服务发现所有节点
			service, err := discovery.GetService(ctx, clientMeta.ServiceName)
			if err != nil {
				return nil,rpcConst.NotFoundNode
			}
			logBase.Debug("DiscoveryMiddleware,service=%v",service)
			//服务路由：根据模式找到符合要求的节点集合
			allNode := getAllNodes(clientMeta,service)
			logBase.Debug("DiscoveryMiddleware,allNode=%v",allNode)
			//保存节点
			clientMeta.AllNodes = allNode
			return next(ctx, req)
		}
	}
}

func getAllNodes(clientMeta *meta.ClientMeta,service *register.Service) (nodes []*register.Node){
	allNode := service.SvrNodes
	switch clientMeta.CallerType{
	case meta.Caller_type_balance:
		return allNode
	case meta.Caller_type_one:
		svrId := clientMeta.CallerServerId
		for _,v := range allNode{
			if v.NodeId == svrId {
				return []*register.Node{v}
			}
		}
		return nil
	case meta.Caller_type_all:
		return allNode
	}
	return allNode
}
