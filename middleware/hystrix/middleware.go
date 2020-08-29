package mwHystrix

import (
	"context"
	"github.com/afex/hystrix-go/hystrix"
	"myRPC/log/base"
	"myRPC/meta"
	"myRPC/middleware/base"
)

func HystrixMiddleware() mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			clientMeta := meta.GetClientMeta(ctx)
			//无法连接，熔断
			hystrixErr := hystrix.Do(clientMeta.ServiceName, func() (err error) {
				resp, err = next(ctx, req)
				return err
			}, nil)
			if hystrixErr != nil {
				logBase.Debug("HystrixMiddleware,ServiceName=%s,stop")
				return nil, hystrixErr
			}
			logBase.Debug("HystrixMiddleware,ServiceName=%s,continue")
			return resp, hystrixErr
		}
	}
}