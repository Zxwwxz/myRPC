package limitBase

import (
	"errors"
	"myRPC/limit/limiter"
)

var limitManager *LimitManager

type LimitManager struct {}

//启动服务时初始化
func InitLimit()  {
	limitManager = &LimitManager{}
}

//获取全局限流管理器
func GetLimitMgr()*LimitManager {
	return limitManager
}

//每次rpc调用创建一个限流器
func (l *LimitManager)NewLimiter(limiterType string,params map[interface{}]interface{}) (limiter.LimitInterface,error) {
	switch limiterType {
	case "token":
		return limiter.NewTokenLimit(params),nil
	case "counter":
		return limiter.NewCounterLimit(params),nil
	case "bucket":
		return limiter.NewBucketLimit(params),nil
	}
	return nil,errors.New("limiterType illegal")
}
