package mwLimit

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"myRPC/const"
	"myRPC/limit/limiter"
	"myRPC/log/base"
	"myRPC/middleware/base"
)

//客户端限流中间件
func ClientLimitMiddleware(limiter limiter.LimitInterface) mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			//判断是否有限流
			allow := limiter.Allow()
			logBase.Debug("ClientLimitMiddleware,allow=%s",allow)
			if !allow {
				return nil,rpcConst.ClientLimit
			}
			return next(ctx, req)
		}
	}
}

//服务端限流中间件
func ServerLimitMiddleware(limiter limiter.LimitInterface) mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			//判断是否有限流
			allow := limiter.Allow()
			logBase.Debug("ServerLimitMiddleware,allow=%s",allow)
			if !allow {
				return nil,status.Error(codes.ResourceExhausted, "server rate limited")
			}
			return next(ctx, req)
		}
	}
}
