package mwLimit

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"myRPC/limit/limiter"
	logBase "myRPC/log/base"
	mwBase "myRPC/middleware/base"
)

func LimitMiddleware(limiter limiter.LimitInterface) mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			//判断是否有限流
			allow := limiter.Allow()
			logBase.Debug("LimitMiddleware,alloc=%s",allow)
			if !allow {
				err = status.Error(codes.ResourceExhausted, "rate limited")
				return
			}
			return next(ctx, req)
		}
	}
}
