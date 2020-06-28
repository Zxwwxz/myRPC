package mwDiscover

import (
	"context"
	"fmt"
	"myRPC/meta"
	mwBase "myRPC/middleware/base"
	registryBase "myRPC/registry/base"
)

func DiscoveryMiddleware(discovery registryBase.RegistryPlugin) mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			clientMeta := meta.GetClientMeta(ctx)
			//服务发现所有节点
			service, err := discovery.GetService(ctx, clientMeta.ServiceName)
			if err != nil {
				return
			}
			clientMeta.AllNodes = service.SvrNodes
			fmt.Println("进入服务发现中间件：",clientMeta.AllNodes)
			resp, err = next(ctx, req)
			return
		}
	}
}
