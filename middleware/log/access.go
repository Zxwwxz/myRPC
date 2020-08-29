package mwLog

import (
	"context"
	"google.golang.org/grpc/status"
	logBase "myRPC/log/base"
	"myRPC/meta"
	mwBase "myRPC/middleware/base"
	"time"
)

func LogServiceMiddleware() mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			startTime := time.Now()
			serverMeta := meta.GetServerMeta(ctx)
			logBase.Debug("LogServiceMiddleware->befor serverMeta:%v",serverMeta)
			resp, err = next(ctx, req)
			errStatus, _ := status.FromError(err)
			cost := time.Since(startTime).Nanoseconds() / 1000
			logBase.Debug("LogServiceMiddleware->after serverMeta:%v",serverMeta)
			logBase.Debug("LogServiceMiddleware,code:%d",errStatus.Code())
			logBase.Debug("LogServiceMiddleware,cost:%d",cost)
			return
		}
	}
}
func LogClientMiddleware() mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			startTime := time.Now()
			clientMeta := meta.GetClientMeta(ctx)
			logBase.Debug("LogClientMiddleware->befor clientMeta:%v",clientMeta)
			resp, err = next(ctx, req)
			errStatus, _ := status.FromError(err)
			cost := time.Since(startTime).Nanoseconds() / 1000
			logBase.Debug("LogClientMiddleware->after clientMeta:%v",clientMeta)
			logBase.Debug("LogClientMiddleware,code:%d",errStatus.Code())
			logBase.Debug("LogClientMiddleware,cost:%d",cost)
			return
		}
	}
}

