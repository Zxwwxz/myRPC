package mwLog

import (
	"context"
	"myRPC/log/base"
	"myRPC/meta"
	"myRPC/middleware/base"
	"time"
)

func LogServiceMiddleware() mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			startTime := time.Now()
			serverMeta := meta.GetServerMeta(ctx)
			logBase.Debug("LogServiceMiddleware->befor serverMeta:%v",serverMeta)
			resp, err = next(ctx, req)
			cost := time.Since(startTime).Nanoseconds() / 1000
			logBase.Debug("LogServiceMiddleware->after serverMeta:%v",serverMeta)
			if err != nil {
				logBase.Debug("LogServiceMiddleware,errMsg:%s",err.Error())
			}
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
			cost := time.Since(startTime).Nanoseconds() / 1000
			logBase.Debug("LogClientMiddleware->after clientMeta:%v",clientMeta)
			if err != nil {
				logBase.Debug("LogClientMiddleware,errMsg:%s",err.Error())
			}
			logBase.Debug("LogClientMiddleware,cost:%d",cost)
			return
		}
	}
}

