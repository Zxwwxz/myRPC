package mwPrometheus

import (
	"context"
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
			DefaultServerMetrics.IncRequest(ctx, serverMeta.ServiceName, serverMeta.ServiceMethod)
			startTime := time.Now()
			resp, err = next(ctx, req)
			DefaultServerMetrics.IncErrcode(ctx, serverMeta.ServiceName, serverMeta.ServiceMethod, err)
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
			DefaultClientMetrics.IncRequest(ctx, clientMeta.ServiceName, clientMeta.ServiceMethod)
			startTime := time.Now()
			resp, err = next(ctx, req)
			DefaultClientMetrics.IncErrcode(ctx, clientMeta.ServiceName, clientMeta.ServiceMethod, err)
			DefaultClientMetrics.ObserveLatency(ctx, clientMeta.ServiceName,
				clientMeta.ServiceMethod, time.Since(startTime).Nanoseconds()/1000)
			return
		}
	}
}