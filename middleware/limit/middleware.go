package mwLimit

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"myRPC/limit"
	mwBase "myRPC/middleware/base"
)

func LimitMiddleware(limiter limit.LimitInterface) mwBase.MiddleWare {
	return func(next mwBase.MiddleWareFunc) mwBase.MiddleWareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			//判断是否有限流
			allow := limiter.Allow()
			fmt.Println("进入限流中间件：",allow)
			if !allow {
				err = status.Error(codes.ResourceExhausted, "rate limited")
				return
			}
			return next(ctx, req)
		}
	}
}
