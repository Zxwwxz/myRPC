package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"myRPC/hystrix"
	balanceBase "myRPC/loadBalance/base"
	logBase "myRPC/log/base"
	limitBase "myRPC/limit/base"
	registryBase "myRPC/registry/base"
	"myRPC/prometheus"
	"myRPC/config"
	"myRPC/limit/limiter"
	"myRPC/registry/register"
	"myRPC/trace"
	mwTrace "myRPC/middleware/trace"
	mwBase "myRPC/middleware/base"
	mwLimit "myRPC/middleware/limit"
	mwLog "myRPC/middleware/log"
	mwPrometheus "myRPC/middleware/prometheus"
	"myRPC/meta"
	"myRPC/util"
	"net"
)

var commonService *CommonService

type CommonService struct {
	*grpc.Server
	serviceConf *config.ServiceConf
	limiter     limiter.LimitInterface
}

func InitService()(err error) {
	//创建公共服务对象
	commonService = &CommonService{
		Server:grpc.NewServer(),
	}
	//初始化配置
	err = config.InitConfig()
	if err != nil {
		return
	}
	commonService.serviceConf = config.GetConf()
	//初始化
	err = initLimit()
	if err != nil {
		return
	}
	err = initLog()
	if err != nil {
		return
	}
	err = initRegistry()
	if err != nil {
		return
	}
	err = initTrace()
	if err != nil {
		return
	}
	err = initPrometheus()
	if err != nil {
		return
	}
	err = initBalance()
	if err != nil {
		return
	}
	err = initHystrix()
	if err != nil {
		return
	}
	return
}

func InitServiceFunc(reqCtx context.Context,serviceMethod string)(ctx context.Context,err error) {
	serverMeta := &meta.ServerMeta{
		Env:util.GetEnv(),
		IDC:commonService.serviceConf.Base.ServiceIDC,
		ServeiceIP:util.GetLocalIP(),
		ServiceName:commonService.serviceConf.Base.ServiceName,
		ServiceMethod:serviceMethod,
	}
	ctx = meta.SetServerMeta(reqCtx,serverMeta)
	return ctx,nil
}

func Run() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", commonService.serviceConf.Base.ServicePort))
	if err != nil {
		fmt.Println("listen err:",err)
	}
	err = commonService.Server.Serve(listen)
	if err != nil {
		fmt.Println("start service err:",err)
	}
}
func GetGrpcService() *grpc.Server {
	return commonService.Server
}

func initLimit()(err error) {
	limitBase.InitLimit()
	if commonService.serviceConf.ServerLimit.SwitchOn == true{
		serverLimiter,err := limitBase.GetLimitMgr().NewLimiter(commonService.serviceConf.ServerLimit.Type,
			commonService.serviceConf.ServerLimit.Params.(map[interface{}]interface{}))
		limitBase.GetLimitMgr().SetServerLimiter(serverLimiter)
		commonService.limiter = serverLimiter
		if err != nil {
			return err
		}
	}
	if commonService.serviceConf.ClientLimit.SwitchOn == true{
		clientLimiter,err := limitBase.GetLimitMgr().NewLimiter(commonService.serviceConf.ClientLimit.Type,
			commonService.serviceConf.ClientLimit.Params.(map[interface{}]interface{}))
		limitBase.GetLimitMgr().SetClientLimiter(clientLimiter)
		if err != nil {
			return err
		}
	}
	return err
}

func initLog()(err error) {
	if commonService.serviceConf.Log.SwitchOn == false{
		return nil
	}
	logBase.InitLogger(commonService.serviceConf.Log.Level,
		commonService.serviceConf.Log.ChanSize,
		commonService.serviceConf.Log.Params.(map[interface{}]interface{}))
	return
}

func initRegistry()(err error) {
	registryBase.InitRegistry()
	_,err = registryBase.GetRegistryManager().NewRegister(commonService.serviceConf.Registry.Type,
		commonService.serviceConf.Registry.Params.(map[interface{}]interface{}))
	if err != nil {
		return err
	}
	registerServer := &register.Service{
		SvrName:commonService.serviceConf.Base.ServiceName,
		SvrType:commonService.serviceConf.Base.ServiceType,
		SvrNodes:[]*register.Node{
			&register.Node{
				NodeIDC:commonService.serviceConf.Base.ServiceIDC,
				NodeId:commonService.serviceConf.Base.ServiceId,
				NodeVersion:commonService.serviceConf.Base.ServiceVer,
				NodeIp:util.GetLocalIP(),
				NodePort:fmt.Sprintf("%d",commonService.serviceConf.Base.ServicePort),
				NodeWeight:commonService.serviceConf.Base.ServiceWidget,
				NodeFuncs:commonService.serviceConf.Base.ServiceFuncs,
			},
		},
	}
	err = registryBase.GetRegistryManager().RegisterServer(registerServer)
	return err
}

func initTrace()(err error) {
	if commonService.serviceConf.Trace.SwitchOn == false{
		return nil
	}
	err = trace.InitTrace(commonService.serviceConf.Base.ServiceName,
		commonService.serviceConf.Trace.ReportAddr,
		commonService.serviceConf.Trace.SampleType,
		commonService.serviceConf.Trace.SampleRate)
	if err != nil {
		return err
	}
	return
}

func initPrometheus()(err error) {
	if commonService.serviceConf.Prometheus.SwitchOn {
		return prometheus.InitPrometheus(commonService.serviceConf.Prometheus.ListenPort,
			commonService.serviceConf.Prometheus.ClientHistogram,
			commonService.serviceConf.Prometheus.ServerHistogram)
	}
	return nil
}

func initBalance()(error) {
	balanceBase.InitBalance()
	_,err := balanceBase.GetBalanceMgr().NewBalancer(commonService.serviceConf.Balance.Type)
	return err
}

func initHystrix()(error)  {
	err := hystrix.InitHystrix(commonService.serviceConf.Base.ServiceName,
		commonService.serviceConf.Hystrix.TimeOut,
		commonService.serviceConf.Hystrix.MaxConcurrentRequests,
		commonService.serviceConf.Hystrix.RequestVolumeThreshold,
		commonService.serviceConf.Hystrix.SleepWindow,
		commonService.serviceConf.Hystrix.ErrorPercentThreshold,
		)
	return err
}

//服务中间件
func BuildServerMiddleware(handle mwBase.MiddleWareFunc,frontMiddles,backMiddles []mwBase.MiddleWare) mwBase.MiddleWareFunc {
	var middles []mwBase.MiddleWare
	serviceConf := config.GetConf()
	middles = append(middles,frontMiddles...)
	middles = append(middles, mwLog.LogServiceMiddleware())
	if serviceConf.Prometheus.SwitchOn {
		middles = append(middles, mwPrometheus.PrometheusServiceMiddleware())
	}
	if serviceConf.ServerLimit.SwitchOn && commonService.limiter != nil{
		middles = append(middles, mwLimit.ServerLimitMiddleware(commonService.limiter))
	}
	if serviceConf.Trace.SwitchOn {
		middles = append(middles, mwTrace.TraceServiceMiddleware())
	}
	middles = append(middles,backMiddles...)
	m := mwBase.Chain(middles...)
	return m(handle)
}

func Stop() {

}