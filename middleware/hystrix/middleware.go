package mwHystrix

import (
	"context"
	"fmt"
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
			hystrixErr := hystrix.Do(clientMeta.ClientName, func() (err error) {
				resp, err = next(ctx, req)
				return err
			}, nil)
			circuit,_,_ := hystrix.GetCircuit(clientMeta.ClientName)
			fmt.Printf("HystrixMiddleware,ServiceName=%s,isOpen=%v,isAllow=%v",clientMeta.ClientName,circuit.IsOpen(),circuit.AllowRequest())
			fmt.Println()
			logBase.Debug("HystrixMiddleware,ServiceName=%s,isOpen=%v,isAllow=%v",clientMeta.ClientName,circuit.IsOpen(),circuit.AllowRequest())
			if hystrixErr != nil {
				logBase.Debug("HystrixMiddleware,ServiceName=%s,err:%s",clientMeta.ClientName,hystrixErr.Error())
				return nil, hystrixErr
			}
			return resp, nil
		}
	}
}