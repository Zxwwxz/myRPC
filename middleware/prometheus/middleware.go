package mwPrometheus

import (
	"context"
	"fmt"
	"myRPC/meta"
	mwBase "myRPC/middleware/base"
	"time"
)

var (
	DefaultServerMetrics = NewServerMetrics()
	DefaultClientMetrics = NewClientMetrics()
)

func PrometheusServiceMiddleware() mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			serverMeta := meta.GetServerMeta(ctx)
			//监控调用量
			DefaultServerMetrics.IncRequest(ctx, serverMeta.ServiceName, serverMeta.ServiceMethod)
			startTime := time.Now()
			resp, err = next(ctx, req)
			//监控错误码
			DefaultServerMetrics.IncErrcode(ctx, serverMeta.ServiceName, serverMeta.ServiceMethod, err)
			//耗时分布
			DefaultServerMetrics.ObserveLatency(ctx, serverMeta.ServiceName,
				serverMeta.ServiceMethod, time.Since(startTime).Nanoseconds()/1000)
			return
		}
	}
}

func PrometheusClientMiddleware() mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			clientMeta := meta.GetClientMeta(ctx)
			//监控调用量
			DefaultClientMetrics.IncRequest(ctx, clientMeta.ServiceName, clientMeta.ServiceMethod)
			fmt.Println("进入普罗米修斯中间件：",clientMeta.ServiceName,clientMeta.ServiceMethod)
			startTime := time.Now()
			resp, err = next(ctx, req)
			//监控错误码
			DefaultClientMetrics.IncErrcode(ctx, clientMeta.ServiceName, clientMeta.ServiceMethod, err)
			//耗时分布
			DefaultClientMetrics.ObserveLatency(ctx, clientMeta.ServiceName,
				clientMeta.ServiceMethod, time.Since(startTime).Nanoseconds()/1000)
			return
		}
	}
}