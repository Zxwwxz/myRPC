package mwHystrix

import (
	"context"
	"github.com/afex/hystrix-go/hystrix"
	"myRPC/log/base"
	"myRPC/meta"
	"myRPC/middleware/base"
)

//熔断中间件
func HystrixMiddleware() mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			clientMeta := meta.GetClientMeta(ctx)
			hystrixErr := hystrix.Do(clientMeta.ClientName, func() (err error) {
				//没问题
				resp, err = next(ctx, req)
				return err
			}, nil)
			//无法连接，获取熔断状态
			circuit,_,_ := hystrix.GetCircuit(clientMeta.ClientName)
			logBase.Debug("HystrixMiddleware,ServiceName=%s,isOpen=%v,isAllow=%v",clientMeta.ClientName,circuit.IsOpen(),circuit.AllowRequest())
			if hystrixErr != nil {
				logBase.Debug("HystrixMiddleware,ServiceName=%s,err:%s",clientMeta.ClientName,hystrixErr.Error())
				return nil, hystrixErr
			}
			return resp, nil
		}
	}
}