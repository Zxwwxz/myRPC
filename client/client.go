package client

import (
	"context"
	"myRPC/config"
	limitBase "myRPC/limit/base"
	"myRPC/limit/limiter"
	"myRPC/loadBalance/balancer"
	balanceBase "myRPC/loadBalance/base"
	mwBase "myRPC/middleware/base"
	wmConn "myRPC/middleware/conn"
	mwDiscover "myRPC/middleware/discover"
	mwHystrix "myRPC/middleware/hystrix"
	mwLimit "myRPC/middleware/limit"
	mwLoadBalance "myRPC/middleware/loadBalance"
	mwLog "myRPC/middleware/log"
	mwPrometheus "myRPC/middleware/prometheus"
	mwTrace "myRPC/middleware/trace"
	registryBase "myRPC/registry/base"
	"myRPC/registry/register"
)

//公共客户端调用
type CommonClient struct {
	//服务配置
	serviceConf *config.ServiceConf
	//服务限流
	limiter  limiter.LimitInterface
	//服务负载
	balancer  balancer.BalanceInterface
	//服务注册
	register  register.RegisterInterface
}

func InitClient() (*CommonClient,error) {
	//初始配置
	commonClient := &CommonClient{}
	commonClient.serviceConf = config.GetConf()
	err := commonClient.initLimit()
	if err != nil {
		return nil,err
	}
	err = commonClient.initBalance()
	if err != nil {
		return nil,err
	}
	err = commonClient.initRegistry()
	if err != nil {
		return nil,err
	}
	return commonClient,nil
}

func InitClientFunc(reqCtx context.Context) (ctx context.Context,err error) {

}

func (commonClient *CommonClient)initLimit()(error)  {
	if commonClient.serviceConf.Limit.SwitchOn == false{
		return nil
	}
	tempLimiter,err := limitBase.GetLimitMgr().NewLimiter(commonClient.serviceConf.Limit.Type,
		commonClient.serviceConf.Limit.Params.(map[interface{}]interface{}))
	commonClient.limiter = tempLimiter
	return err
}

func (commonClient *CommonClient)initRegistry()(error)  {
	commonClient.register = registryBase.GetRegister()
	return nil
}

func (commonClient *CommonClient)initBalance()(error) {
	tempBalancer,err := balanceBase.GetBalanceMgr().NewBalancer(commonClient.serviceConf.Balance.Type)
	commonClient.balancer = tempBalancer
	return err
}

func (commonClient *CommonClient)BuildClientMiddleware(handle mwBase.MiddleWareFunc,frontMiddles,backMiddles []mwBase.MiddleWare) mwBase.MiddleWareFunc {
	var middles []mwBase.MiddleWare
	//前置中间件
	middles = append(middles,frontMiddles...)
	if commonClient.serviceConf.Log.SwitchOn {
		//日志中间件
		middles = append(middles,mwLog.LogClientMiddleware())
	}
	if commonClient.serviceConf.Limit.SwitchOn {
		//限流中间件
		middles = append(middles,mwLimit.LimitMiddleware(commonClient.limiter))
	}
	if commonClient.serviceConf.Hystrix.SwitchOn {
		//熔断中间件
		middles = append(middles,mwHystrix.HystrixMiddleware())
	}
	if commonClient.serviceConf.Prometheus.SwitchOn {
		//监控中间件
		middles = append(middles,mwPrometheus.PrometheusClientMiddleware())
	}
	if commonClient.serviceConf.Trace.SwitchOn {
		//追踪id中间件
		middles = append(middles,mwTrace.TraceIdClientMiddleware())
		//追踪中间件
		middles = append(middles,mwTrace.TraceClientMiddleware())
	}
	//服务发现中间件
	middles = append(middles,mwDiscover.DiscoveryMiddleware(commonClient.register))
	//负载均衡中间件
	middles = append(middles,mwLoadBalance.LoadBalanceMiddleware(commonClient.balancer))
	//连接中间件
	middles = append(middles,wmConn.ConnMiddleware())
	//后续中间件
	middles = append(middles,backMiddles...)
	//中间件串联
	m := mwBase.Chain(middles...)
	return m(handle)
}
