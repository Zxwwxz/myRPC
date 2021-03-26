package limitBase

import (
	"errors"
	"myRPC/limit/limiter"
)

var limitManager *LimitManager

type LimitManager struct {
	clientLimiter limiter.LimitInterface
	serverLimiter limiter.LimitInterface
}

//启动服务时初始化
func InitLimit()  {
	limitManager = &LimitManager{}
}

//获取全局限流管理器
func GetLimitMgr()*LimitManager {
	return limitManager
}

//客户端服务器限流对象
func (l *LimitManager)SetClientLimiter(clientLimiter limiter.LimitInterface) {
	l.clientLimiter = clientLimiter
}

func (l *LimitManager)GetClientLimiter()(limiter.LimitInterface) {
	return l.clientLimiter
}

func (l *LimitManager)SetServerLimiter(serverLimiter limiter.LimitInterface) {
	l.serverLimiter = serverLimiter
}

func (l *LimitManager)GetServerLimiter()(limiter.LimitInterface) {
	return l.serverLimiter
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

func (l *LimitManager)Stop() () {
	l.serverLimiter = nil
	l.clientLimiter = nil
}
