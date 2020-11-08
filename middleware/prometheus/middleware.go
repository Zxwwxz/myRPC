package mwPrometheus

import (
	"context"
	"google.golang.org/grpc/status"
	"myRPC/const"
	"myRPC/log/base"
	"myRPC/meta"
	"myRPC/middleware/base"
	"myRPC/prometheus"
	"time"
)

func PrometheusServiceMiddleware() mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			serverMeta := meta.GetServerMeta(ctx)
			//监控调用量
			logBase.Debug("PrometheusServiceMiddleware")
			serverMetrics := prometheus.GetPrometheusManager().GetServerMetrics()
			serverMetrics.IncRequest(ctx, serverMeta.ServiceName, serverMeta.ServiceMethod)
			startTime := time.Now()
			resp, err = next(ctx, req)
			if err != nil {
				//监控错误码
				serverMetrics.IncErrcode(ctx, serverMeta.ServiceName, serverMeta.ServiceMethod,"server", int(status.Code(err)))
			}
			//耗时分布
			serverMetrics.ObserveLatency(ctx, serverMeta.ServiceName,
				serverMeta.ServiceMethod, time.Since(startTime).Nanoseconds()/1000000)
			return
		}
	}
}

func PrometheusClientMiddleware() mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			clientMeta := meta.GetClientMeta(ctx)
			logBase.Debug("PrometheusClientMiddleware")
			clientMetrics := prometheus.GetPrometheusManager().GetClientMetrics()
			//监控调用量
			clientMetrics.IncRequest(ctx, clientMeta.ServiceName, clientMeta.ServiceMethod)
			startTime := time.Now()
			resp, err = next(ctx, req)
			if newErr,ok := err.(*rpcConst.ClientError);ok == true{
				//监控错误码
				clientMetrics.IncErrcode(ctx, clientMeta.ServiceName, clientMeta.ServiceMethod, "client", newErr.Code())
			}
			//耗时分布
			clientMetrics.ObserveLatency(ctx, clientMeta.ServiceName,
				clientMeta.ServiceMethod, time.Since(startTime).Nanoseconds()/1000)
			return
		}
	}
}