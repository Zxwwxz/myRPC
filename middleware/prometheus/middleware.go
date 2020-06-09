package mwPrometheus

import (
	"context"
	"myRPC/meta"
	mwBase "myRPC/middleware/base"
	"time"
)

var (
	DefaultServerMetrics = NewServerMetrics()
)

func PrometheusServerMiddleware(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
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
