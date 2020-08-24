package mwLog

import (
	"context"
	"fmt"
	"google.golang.org/grpc/status"
	logBase "myRPC/log/base"
	logOutputer "myRPC/log/outputer"
	"myRPC/meta"
	mwBase "myRPC/middleware/base"
	"time"
)

func LogServiceMiddleware() mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			startTime := time.Now()
			serverMeta := meta.GetServerMeta(ctx)
			logBase.Debug("LogServiceMiddleware->serverMeta:%v",serverMeta)
			resp, err = next(ctx, req)
			errStatus, _ := status.FromError(err)
			cost := time.Since(startTime).Nanoseconds() / 1000
			logBase.Debug("LogServiceMiddleware,code:%d",errStatus.Code())
			logBase.Debug("LogServiceMiddleware,cost:%d",cost)
			return
		}
	}
}
func LogClientMiddleware() mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			ctx = logBase.WithFieldContext(ctx)
			startTime := time.Now()
			fmt.Println("进入日志中间件：",startTime)
			resp, err = next(ctx, req)
			clientMeta := meta.GetClientMeta(ctx)
			errStatus, _ := status.FromError(err)
			cost := time.Since(startTime).Nanoseconds() / 1000
			logBase.AddField(ctx, "client:cost", cost)
			logBase.AddField(ctx, "client:service_method", clientMeta.ServiceMethod)
			logBase.AddField(ctx, "client:server_name", clientMeta.ServiceName)
			logBase.AddField(ctx, "client:caller_cluster", clientMeta.CallerCluster)
			logBase.AddField(ctx, "client:service_cluster", clientMeta.ServiceCluster)
			logBase.AddField(ctx, "client:env", clientMeta.Env)
			if clientMeta.curNode != nil {
				logBase.AddField(ctx, "client:select_node", fmt.Sprintf("%s:%s,", clientMeta.CurNode.NodeIp, clientMeta.CurNode.NodePort))
			}
			logBase.AddField(ctx, "client:caller_idc", clientMeta.CallerIDC)
			logBase.AddField(ctx, "client:service_idc", clientMeta.ServiceIDC)
			logOutputer.Access(ctx, "client:result=%v", errStatus.Code())
			return
		}
	}
}

