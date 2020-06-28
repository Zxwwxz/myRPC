package mwHystrix

import (
	"context"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"myRPC/meta"
	mwBase "myRPC/middleware/base"
)

func HystrixMiddleware() mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			clientMeta := meta.GetClientMeta(ctx)
			//无法连接，熔断
			fmt.Println("进入熔断中间件：",clientMeta.ServiceName)
			hystrixErr := hystrix.Do(clientMeta.ServiceName, func() (err error) {
				resp, err = next(ctx, req)
				return err
			}, nil)
			if hystrixErr != nil {
				return nil, hystrixErr
			}
			return resp, hystrixErr
		}
	}
}