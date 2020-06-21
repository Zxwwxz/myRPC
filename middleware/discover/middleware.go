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
			rpcMeta := meta.GetClientMeta(ctx)
			rpcMeta.AllNodes = make(map[int]*registryBase.Node)
			service, err := discovery.GetService(ctx, rpcMeta.ServiceName)
			if err != nil {
				return
			}
			rpcMeta.AllNodes = service.SvrNodes
			resp, err = next(ctx, req)
			return
		}
	}
}
