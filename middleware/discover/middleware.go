package mwDiscover

import (
	"context"
	"myRPC/meta"
	mwBase "myRPC/middleware/base"
	registryBase "myRPC/registry/base"
)

func DiscoveryMiddleware(discovery registryBase.RegistryPlugin) mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			clientMeta := meta.GetClientMeta(ctx)
			clientMeta.AllNodes = make(map[int]*registryBase.Node)
			//服务发现所有节点
			service, err := discovery.GetService(ctx, clientMeta.ServiceName)
			if err != nil {
				return
			}
			clientMeta.AllNodes = service.SvrNodes
			resp, err = next(ctx, req)
			return
		}
	}
}
