package service

import (
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"myRPC/http"
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
	"os"
	"os/signal"
)

//基础服务模块，所有服务的公共部分
type CommonService struct {
	//grpc服务对象
	*grpc.Server
	//服务配置对象
	serviceConf *config.ServiceConf
	//服务端限流器
	limiter     limiter.LimitInterface
	//http服务器
	httpServer  *httpBase.HttpServer
}

//新建服务对象
func NewService()(commonService *CommonService,err error) {
	//创建公共服务对象
	commonService = &CommonService{
		Server:grpc.NewServer(),
	}
	//初始化命令行参数
	configDir,serviceParams,err := commonService.initParams()
	if err != nil {
		return
	}
	//初始化配置
	commonService.serviceConf,err = config.NewConfig(configDir,serviceParams)
	if err != nil {
		return
	}
	fmt.Println("service serviceConf:",commonService.serviceConf)
	//初始化各类对象
	err = commonService.initHttp()
	if err != nil {
		return
	}
	err = commonService.initLimit()
	if err != nil {
		return
	}
	err = commonService.initLog()
	if err != nil {
		return
	}
	err = commonService.initRegistry()
	if err != nil {
		return
	}
	err = commonService.initTrace()
	if err != nil {
		return
	}
	err = commonService.initPrometheus()
	if err != nil {
		return
	}
	err = commonService.initBalance()
	if err != nil {
		return
	}
	err = commonService.initHystrix()
	if err != nil {
		return
	}
	return commonService,nil
}

//服务的接口初始化
func (commonService *CommonService)InitServiceFunc(reqCtx context.Context,serviceMethod string)(ctx context.Context,err error) {
	//服务端执行中间件参数
	serverMeta := &meta.ServerMeta{
		Env:util.GetEnv(),
		IDC:commonService.serviceConf.Base.ServiceIDC,
		ServeiceIP:util.GetLocalIP(),
		ServiceName:commonService.serviceConf.Base.ServiceName,
		ServiceMethod:serviceMethod,
	}
	//保存到上下文中
	ctx = meta.SetServerMeta(reqCtx,serverMeta)
	return ctx,nil
}

//执行服务
func (commonService *CommonService)Run() {
	logBase.Debug("init server start")
	if commonService.httpServer != nil {
		go func() {
			//开启http服务
			err := commonService.httpServer.Start()
			if err != nil {
				logBase.Fatal("start http err:%v",err)
				return
			}
		}()
	}
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", commonService.serviceConf.Base.ServicePort))
	if err != nil {
		logBase.Fatal("new listen err:%v",err)
		return
	}
	//开启rpc服务
	err = commonService.Server.Serve(listen)
	if err != nil {
		logBase.Fatal("start server err:%v",err)
		return
	}
}

//解析命令行可能传递的参数
func (commonService *CommonService)initParams()(configDir string,serviceParams config.ServiceParams,err error)  {
	serviceParams = config.ServiceParams{}
	flag.StringVar(&configDir,"config","","service config path")
	flag.IntVar(&serviceParams.ServiceType,"type",0,"service type")
	flag.IntVar(&serviceParams.ServiceId,"id",0,"service id")
	flag.IntVar(&serviceParams.ServiceVer,"ver",0,"service ver")
	flag.StringVar(&serviceParams.ServiceName,"name","","service name")
	flag.IntVar(&serviceParams.ServicePort,"sport",0,"service port")
	flag.IntVar(&serviceParams.HttpPort,"hport",0,"service http port")
	flag.Parse()
	return "",serviceParams,nil
}

//初始化限流
func (commonService *CommonService)initLimit()(err error) {
	limitBase.InitLimit()
	//服务端限流器
	serverLimiter,err := limitBase.GetLimitMgr().NewLimiter(commonService.serviceConf.ServerLimit.Type,
		commonService.serviceConf.ServerLimit.Params.(map[interface{}]interface{}))
	limitBase.GetLimitMgr().SetServerLimiter(serverLimiter)
	commonService.limiter = serverLimiter
	if err != nil {
		return err
	}
	//客户端限流器
	clientLimiter,err := limitBase.GetLimitMgr().NewLimiter(commonService.serviceConf.ClientLimit.Type,
		commonService.serviceConf.ClientLimit.Params.(map[interface{}]interface{}))
	limitBase.GetLimitMgr().SetClientLimiter(clientLimiter)
	if err != nil {
		return err
	}
	return nil
}

//初始化日志
func (commonService *CommonService)initLog()(err error) {
	logBase.InitLogger(commonService.serviceConf.Log.SwitchOn,
		commonService.serviceConf.Log.Level,
		commonService.serviceConf.Log.ChanSize,
		commonService.serviceConf.Log.Params.(map[interface{}]interface{}))
	return
}

//初始化注册中心
func (commonService *CommonService)initRegistry()(err error) {
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
	//起服的时候，注册当前服务到注册中心
	err = registryBase.GetRegistryManager().RegisterServer(registerServer)
	return err
}

//初始化追踪
func (commonService *CommonService)initTrace()(err error) {
	err = trace.InitTrace(commonService.serviceConf.Base.ServiceName,
		commonService.serviceConf.Trace.ReportAddr,
		commonService.serviceConf.Trace.SampleType,
		commonService.serviceConf.Trace.SampleRate)
	if err != nil {
		return err
	}
	return
}

//初始化监控
func (commonService *CommonService)initPrometheus()(err error) {
	err = prometheus.NewPrometheusManager(
		//客户端监控频率
		commonService.serviceConf.Prometheus.ClientHistogram,
		//服务端监控频率
		commonService.serviceConf.Prometheus.ServerHistogram)
	if err != nil {
		return err
	}
	if commonService.httpServer != nil {
		//添加http接口调用
		err = prometheus.GetPrometheusManager().AddPrometheusHandler(commonService.httpServer.GetRoute(),commonService.serviceConf)
		if err != nil {
			return err
		}
	}
	return nil
}

//初始化负载均衡
func (commonService *CommonService)initBalance()(error) {
	balanceBase.InitBalance()
	_,err := balanceBase.GetBalanceMgr().NewBalancer(commonService.serviceConf.Balance.Type)
	return err
}

//初始化熔断
func (commonService *CommonService)initHystrix()(error)  {
	err := hystrix.InitHystrix(commonService.serviceConf.Base.ServiceName,
		commonService.serviceConf.Hystrix.TimeOut,
		commonService.serviceConf.Hystrix.MaxConcurrentRequests,
		commonService.serviceConf.Hystrix.RequestVolumeThreshold,
		commonService.serviceConf.Hystrix.SleepWindow,
		commonService.serviceConf.Hystrix.ErrorPercentThreshold,
		)
	return err
}

//初始化http服务
func (commonService *CommonService)initHttp()(err error)  {
	if commonService.serviceConf.Http.SwitchOn {
		httpServer,err := httpBase.NewHttpServer(commonService.serviceConf.Http.HttpPort)
		if err != nil {
			return err
		}
		commonService.httpServer = httpServer
		//开启pprof性能监控
		if commonService.serviceConf.Http.PprofSwitchOn {
			err = httpServer.AddPropHandler()
			if err != nil {
				return err
			}
		}
		//添加http配置修改
		err = httpServer.AddParamsHandler(commonService.serviceConf)
		if err != nil {
			return err
		}
	}
	return nil
}

//服务中间件
func (commonService *CommonService)BuildServerMiddleware(handle mwBase.MiddleWareFunc,frontMiddles,backMiddles []mwBase.MiddleWare) mwBase.MiddleWareFunc {
	var middles []mwBase.MiddleWare
	serviceConf := commonService.serviceConf
	//前置
	middles = append(middles,frontMiddles...)
	//日志
	middles = append(middles, mwLog.LogServiceMiddleware())
	//监控
	if serviceConf.Prometheus.SwitchOn {
		middles = append(middles, mwPrometheus.PrometheusServiceMiddleware())
	}
	//服务端限流
	if serviceConf.ServerLimit.SwitchOn && commonService.limiter != nil{
		middles = append(middles, mwLimit.ServerLimitMiddleware(commonService.limiter))
	}
	//追踪
	if serviceConf.Trace.SwitchOn {
		middles = append(middles, mwTrace.TraceServiceMiddleware())
	}
	//后置
	middles = append(middles,backMiddles...)
	m := mwBase.Chain(middles...)
	//得到中间件链条，调用这个返回值会执行所有中间件
	return m(handle)
}

//获取服务配置
func (commonService *CommonService)GetServiceConf()(*config.ServiceConf) {
	return commonService.serviceConf
}

//获取http服务路由
func (commonService *CommonService)GetHttpRouter()(router *mux.Router) {
	if commonService.httpServer != nil {
		return commonService.httpServer.GetRoute()
	}
	return nil
}

//停止服务
func (commonService *CommonService)Stop() {
	stopChan := make(chan os.Signal)
	//监听所有信号
	signal.Notify(stopChan)
	<- stopChan
	logBase.Debug("server stop")
	//各种组件释放资源
	commonService.Server.Stop()
	if commonService.httpServer != nil {
		_ = commonService.httpServer.Stop()
	}
	commonService.serviceConf = nil
	limitBase.GetLimitMgr().Stop()
	registryBase.GetRegistryManager().Stop()
	logBase.GetLogMgr().Stop()
	balanceBase.GetBalanceMgr().Stop()
	prometheus.GetPrometheusManager().Stop()
	hystrix.Stop()
	_ = trace.Stop()
}