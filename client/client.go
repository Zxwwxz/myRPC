package client

import (
	"context"
	"myRPC/config"
	"myRPC/limit/base"
	"myRPC/limit/limiter"
	"myRPC/loadBalance/balancer"
	"myRPC/loadBalance/base"
	"myRPC/meta"
	"myRPC/middleware/base"
	"myRPC/middleware/conn"
	"myRPC/middleware/discover"
	"myRPC/middleware/hystrix"
	"myRPC/middleware/limit"
	"myRPC/middleware/loadBalance"
	"myRPC/middleware/log"
	"myRPC/middleware/prometheus"
	"myRPC/middleware/trace"
	"myRPC/registry/base"
	"myRPC/registry/register"
)

//公共客户端调用
type CommonClient struct {
	//服务配置
	clientConf *config.ServiceConf
	//服务限流
	limiter  limiter.LimitInterface
	//服务负载
	balancer  balancer.BalanceInterface
	//服务注册
	register  register.RegisterInterface
}

func NewClient(reqCtx context.Context,serviceName,serviceMethod string,callerMode int,clientConf *config.ServiceConf,options []meta.ClientMetaOption) (context.Context,*CommonClient,error) {
	//初始配置
	commonClient := &CommonClient{}
	commonClient.clientConf = clientConf
	err := commonClient.initLimit()
	if err != nil {
		return nil,nil,err
	}
	err = commonClient.initBalance()
	if err != nil {
		return nil,nil,err
	}
	err = commonClient.initRegistry()
	if err != nil {
		return nil,nil,err
	}
	ctx,err := commonClient.initClientMeta(reqCtx,serviceName,serviceMethod,callerMode,options)
	if err != nil {
		return nil,nil,err
	}
	return ctx,commonClient,nil
}

func (commonClient *CommonClient)initLimit()(error)  {
	if commonClient.clientConf.ClientLimit.SwitchOn == false{
		return nil
	}
	commonClient.limiter = limitBase.GetLimitMgr().GetClientLimiter()
	return nil
}

func (commonClient *CommonClient)initRegistry()(error)  {
	commonClient.register = registryBase.GetRegister()
	return nil
}

func (commonClient *CommonClient)initBalance()(error) {
	commonClient.balancer = balanceBase.GetCurBalancer()
	return nil
}

func (commonClient *CommonClient)initClientMeta(reqCtx context.Context,serviceName,serviceMethod string,callerMode int,options []meta.ClientMetaOption)(context.Context,error) {
	clientMeta := &meta.ClientMeta{
		ServiceName:serviceName,
		ServiceMethod:serviceMethod,
		ClientName:commonClient.clientConf.Base.ServiceName,
		CallerType:meta.Caller_type_balance,
		CallerMode:callerMode,
		MaxReconnectNum: meta.Default_max_reconnect,
	}
	for _,option := range options{
		option(clientMeta)
	}
	ctx := meta.SetClientMeta(reqCtx,clientMeta)
	return ctx,nil
}

func (commonClient *CommonClient)BuildClientMiddleware(handle mwBase.MiddleWareFunc,frontMiddles,backMiddles []mwBase.MiddleWare) mwBase.MiddleWareFunc {
	var middles []mwBase.MiddleWare
	//前置中间件
	middles = append(middles,frontMiddles...)
	if commonClient.clientConf.Log.SwitchOn {
		//日志中间件
		middles = append(middles,mwLog.LogClientMiddleware())
	}
	if commonClient.clientConf.ClientLimit.SwitchOn && commonClient.limiter != nil{
		//限流中间件
		middles = append(middles,mwLimit.ClientLimitMiddleware(commonClient.limiter))
	}
	if commonClient.clientConf.Hystrix.SwitchOn {
		//熔断中间件
		middles = append(middles,mwHystrix.HystrixMiddleware())
	}
	if commonClient.clientConf.Prometheus.SwitchOn {
		//监控中间件
		middles = append(middles,mwPrometheus.PrometheusClientMiddleware())
	}
	if commonClient.clientConf.Trace.SwitchOn {
		//追踪id中间件
		middles = append(middles,mwTrace.TraceIdClientMiddleware())
		//追踪中间件
		middles = append(middles,mwTrace.TraceClientMiddleware())
	}
	if commonClient.register != nil {
		//服务发现中间件
		middles = append(middles,mwDiscover.DiscoveryMiddleware(commonClient.register))
	}
	if commonClient.balancer != nil {
		//负载均衡中间件
		middles = append(middles,mwLoadBalance.LoadBalanceMiddleware(commonClient.balancer))
	}
	//连接中间件
	middles = append(middles,wmConn.ConnMiddleware())
	//后续中间件
	middles = append(middles,backMiddles...)
	//中间件串联
	m := mwBase.Chain(middles...)
	return m(handle)
}
