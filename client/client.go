package client

import (
	"context"
	"myRPC/config"
	"myRPC/limit"
	"myRPC/loadBalance"
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
	"myRPC/trace"
	"sync"
	"time"
)

var clientOnce sync.Once
var globalRegister registryBase.RegistryPlugin

//公共客户端调用
type CommonClient struct {
	//服务配置
	serviceConf *config.ServiceConf
	//服务注册
	register registryBase.RegistryPlugin
	//服务限流
	limiter  limit.LimitInterface
	//服务负载
	balancer  loadBalance.BalanceInterface
}

func NewCommonClient() (*CommonClient,error) {
	//初始配置
	client := &CommonClient{}
	err := config.InitConfig()
	if err != nil {
		return nil,err
	}
	client.serviceConf = config.GetConf()
	ctx := context.TODO()
	clientOnce.Do(func() {
		//初始全局注册
		regiserConf := client.serviceConf.Regiser
		globalRegister, err = registryBase.PluginManager.InitPlugin(ctx,
			regiserConf.RegisterName,
			registryBase.SetRegisterAddrs([]string{regiserConf.RegisterAddr}),
			registryBase.SetRegisterTimeOut(time.Duration(regiserConf.Timeout)),
			registryBase.SetRegisterPath(regiserConf.RegisterPath),
			registryBase.SetHeartTimeOut(regiserConf.HeartBeat),
		)
		//初始全局追踪
		traceConf := client.serviceConf.Trace
		err = trace.Init(client.serviceConf.ServiceName,traceConf.ReportAddr,traceConf.SampleType,traceConf.SampleRate)
	})
	client.register = globalRegister
	//初始负载
	client.balancer = loadBalance.NewRandomBalance()
	//初始限流
	limitConf := client.serviceConf.Limit
	client.limiter = limit.NewTokenLimit(limitConf.QPSLimit,limitConf.AllWater)
	return client,err
}

func (client *CommonClient)BuildClientMiddleware(handle mwBase.MiddleWareFunc,frontMiddles,backMiddles []mwBase.MiddleWare) mwBase.MiddleWareFunc {
	var middles []mwBase.MiddleWare
	//前置中间件
	middles = append(middles,frontMiddles...)
	//日志中间件
	middles = append(middles,mwLog.AccessClientMiddleware())
	//追踪id中间件
	middles = append(middles,mwTrace.TraceIdClientMiddleware())
	//追踪中间件
	middles = append(middles,mwTrace.TraceClientMiddleware())
	//监控中间件
	middles = append(middles,mwPrometheus.PrometheusClientMiddleware())
	//限流中间件
	middles = append(middles,mwLimit.LimitMiddleware(client.limiter))
	//熔断中间件
	middles = append(middles,mwHystrix.HystrixMiddleware())
	//服务发现中间件
	middles = append(middles,mwDiscover.DiscoveryMiddleware(client.register))
	//负载均衡中间件
	middles = append(middles,mwLoadBalance.LoadBalanceMiddleware(client.balancer))
	//连接中间件
	middles = append(middles,wmConn.ConnMiddleware())
	//后续中间件
	middles = append(middles,backMiddles...)
	//中间件串联
	m := mwBase.Chain(middles...)
	return m(handle)
}
