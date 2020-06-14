package mwLog

import (
	"context"
	"google.golang.org/grpc/status"
	logBase "myRPC/log/base"
	logOutputer "myRPC/log/outputer"
	"myRPC/meta"
	mwBase "myRPC/middleware/base"
	"time"
)

func AccessMiddleware() mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		startTime := time.Now()
		resp, err = next(ctx, req)
		serverMeta := meta.GetServerMeta(ctx)
		errStatus, _ := status.FromError(err)
		cost := time.Since(startTime).Nanoseconds() / 1000
		logBase.AddField(ctx, "cost_us", cost)
		logBase.AddField(ctx, "method", serverMeta.ServiceMethod)
		logBase.AddField(ctx, "cluster", serverMeta.Cluster)
		logBase.AddField(ctx, "env", serverMeta.Env)
		logBase.AddField(ctx, "server_ip", serverMeta.ServerIP)
		logBase.AddField(ctx, "client_ip", serverMeta.ClientIP)
		logBase.AddField(ctx, "idc", serverMeta.IDC)
		logOutputer.Access(ctx, "result=%v", errStatus.Code())
		return
	}
	}
}
