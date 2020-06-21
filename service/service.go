package service

import (
	"context"
	"fmt"
	mwTrace "myRPC/middleware/trace"
	"myRPC/trace"
	"google.golang.org/grpc"
	"myRPC/config"
	"myRPC/limit"
	logOutputer "myRPC/log/outputer"
	mwBase "myRPC/middleware/base"
	mwLimit "myRPC/middleware/limit"
	mwLog "myRPC/middleware/log"
	mwPrometheus "myRPC/middleware/prometheus"
	registryBase "myRPC/registry/base"
	"myRPC/util"
	"net"
	"strconv"
	"time"
)

var commonService = &CommonService{
	Server: grpc.NewServer(),
}

type CommonService struct {
	*grpc.Server
	serviceConf *config.ServiceConf
	Limiter     limit.LimitInterface
}

func Init()  {
	err := config.InitConfig()
	if err != nil {
		fmt.Println("InitConfig,err:",err)
	}
	commonService.serviceConf = config.GetConf()
	initLimit()
	initLogger()
	initRegister()
	initTrace()
}

func Run() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", commonService.serviceConf.Port))
	if err != nil {
		fmt.Println("listen err:",err)
	}
	s := grpc.NewServer()
	commonService.Server = s
	err = s.Serve(listen)
	if err != nil {
		fmt.Println("start service err:",err)
	}
}
func GetGrpcService() *grpc.Server {
	return commonService.Server
}

func initLimit() {
	limiter := limit.NewTokenLimit(commonService.serviceConf.Limit.QPSLimit,commonService.serviceConf.Limit.AllWater)
	commonService.Limiter = limiter
}

func initLogger() {
	filename := fmt.Sprintf("%s/%s.log", commonService.serviceConf.Log.Dir, commonService.serviceConf.ServiceName)
	outputer, err := logOutputer.NewFileOutputer(filename)
	if err != nil {
		return
	}
	logOutputer.InitLogger(commonService.serviceConf.Log.Level, commonService.serviceConf.Log.ChanSize, commonService.serviceConf.ServiceName)
	logOutputer.AddOutputer(outputer)
	return
}

func initRegister() {
	serviceConf := commonService.serviceConf
	registerConf := serviceConf.Regiser
	if !registerConf.SwitchOn {
		return
	}
	localIp := util.GetLocalIP()
	tecdPlugin,err := registryBase.PluginManager.InitPlugin(context.TODO(),
		registerConf.RegisterName,
		registryBase.SetRegisterAddrs([]string{registerConf.RegisterAddr}),
		registryBase.SetRegisterPath(registerConf.RegisterPath),
		registryBase.SetRegisterTimeOut(time.Duration(registerConf.Timeout)*time.Second),
		registryBase.SetHeartTimeOut(registerConf.HeartBeat))
	if err != nil {
		fmt.Println("初始化失败:",err)
		return
	}
	node := &registryBase.Node{NodeId:serviceConf.ServiceId,NodeIp:localIp,NodePort:strconv.Itoa(serviceConf.Port),NodeVersion:serviceConf.ServiceVer,NodeWeight:1,NodeFuncs:[]string{}}
	service := &registryBase.Service{
		SvrName:serviceConf.ServiceName,
		SvrType:serviceConf.ServiceId,
		SvrNodes: map[int]*registryBase.Node{
			serviceConf.ServiceId:node,
		},
	}
	err = tecdPlugin.Register(context.TODO(),service)
	if err != nil {
		fmt.Println("注册失败")
		return
	}
}

func initTrace() {
	serviceConf := commonService.serviceConf
	traceConf := serviceConf.Trace
	if !traceConf.SwitchOn {
		return
	}
	err := trace.Init(serviceConf.ServiceName,traceConf.ReportAddr,traceConf.SampleType,traceConf.SampleRate)
	if err != nil {
		fmt.Println("初始化追踪失败")
		return
	}
}

func BuildServerMiddleware(handle mwBase.MiddleWareFunc,frontMiddles,backMiddles []mwBase.MiddleWare) mwBase.MiddleWareFunc {
	var middles []mwBase.MiddleWare
	serviceConf := config.GetConf()
	middles = append(middles,frontMiddles...)
	middles = append(middles, mwLog.AccessServiceMiddleware())
	if serviceConf.Prometheus.SwitchOn {
		middles = append(middles, mwPrometheus.PrometheusServiceMiddleware())
	}
	if serviceConf.Limit.SwitchOn {
		middles = append(middles, mwLimit.LimitMiddleware(commonService.Limiter))
	}
	if serviceConf.Trace.SwitchOn {
		middles = append(middles, mwTrace.TraceServiceMiddleware())
	}
	middles = append(middles,backMiddles...)
	m := mwBase.Chain(middles...)
	return m(handle)
}
