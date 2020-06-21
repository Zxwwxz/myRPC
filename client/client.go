package client

import (
	"context"
	"fmt"
	"myRPC/config"
	"myRPC/limit"
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
	"myRPC/trace"
	"time"
)

type CommonClient struct {
	serviceConf *config.ServiceConf
	register registryBase.RegistryPlugin
	limiter  limit.LimitInterface
	balancer  balanceBase.BalanceInterface
}

func NewKoalaClient() *CommonClient {
	client := &CommonClient{}
	err := config.InitConfig()
	if err != nil {
		fmt.Println("InitConfig,err:",err)
	}
	client.serviceConf = config.GetConf()
	ctx := context.TODO()
	regiserConf := client.serviceConf.Regiser
	client.register, _ = registryBase.PluginManager.InitPlugin(ctx,
		regiserConf.RegisterName,
		registryBase.SetRegisterAddrs([]string{regiserConf.RegisterAddr}),
		registryBase.SetRegisterTimeOut(time.Duration(regiserConf.Timeout)),
		registryBase.SetRegisterPath(regiserConf.RegisterPath),
		registryBase.SetHeartTimeOut(regiserConf.HeartBeat),
	)
	client.balancer = balanceBase.NewRandomBalance()
	limitConf := client.serviceConf.Limit
	client.limiter = limit.NewTokenLimit(limitConf.QPSLimit,limitConf.AllWater)
	traceConf := client.serviceConf.Trace
	_ = trace.Init(client.serviceConf.ServiceName,traceConf.ReportAddr,traceConf.SampleType,traceConf.SampleRate)
	return client
}

func (client *CommonClient)BuildClientMiddleware(handle mwBase.MiddleWareFunc,frontMiddles,backMiddles []mwBase.MiddleWare) mwBase.MiddleWareFunc {
	var middles []mwBase.MiddleWare
	middles = append(middles,frontMiddles...)
	middles = append(middles,mwTrace.TraceIdClientMiddleware())
	middles = append(middles,mwLog.AccessClientMiddleware())
	middles = append(middles,mwTrace.TraceClientMiddleware())
	middles = append(middles,mwPrometheus.PrometheusClientMiddleware())
	middles = append(middles,mwLimit.LimitMiddleware(client.limiter))
	middles = append(middles,mwHystrix.HystrixMiddleware())
	middles = append(middles,mwDiscover.DiscoveryMiddleware(client.register))
	middles = append(middles,mwLoadBalance.LoadBalanceMiddleware(client.balancer))
	middles = append(middles,wmConn.ConnMiddleware())
	middles = append(middles,backMiddles...)
	m := mwBase.Chain(middles...)
	return m(handle)
}
