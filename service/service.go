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
	"os"
	"strconv"
	"time"
	_ "myRPC/registry/etcd"
)

var commonService = &CommonService{
	Server: grpc.NewServer(),
}

type CommonService struct {
	*grpc.Server
	serviceConf *config.ServiceConf
	Limiter     limit.LimitInterface
}

func Init()(err error)  {
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
	err = initLogger()
	if err != nil {
		return
	}
	err = initRegister()
	if err != nil {
		return
	}
	err = initTrace()
	if err != nil {
		return
	}
	return
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

func initLimit()(err error) {
	limiter := limit.NewTokenLimit(commonService.serviceConf.Limit.QPSLimit,commonService.serviceConf.Limit.AllWater)
	commonService.Limiter = limiter
	return
}

func initLogger()(err error) {
	if !util.IsFileExist(commonService.serviceConf.Log.Dir) {
		err := os.Mkdir(commonService.serviceConf.Log.Dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	filename := fmt.Sprintf("%s/%s.log", commonService.serviceConf.Log.Dir, commonService.serviceConf.ServiceName)
	outputer, err := logOutputer.NewFileOutputer(filename)
	if err != nil {
		return
	}
	logOutputer.InitLogger(commonService.serviceConf.Log.Level, commonService.serviceConf.Log.ChanSize, commonService.serviceConf.ServiceName)
	logOutputer.AddOutputer(outputer)
	return
}

func initRegister()(err error) {
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
		return
	}
	return
}

func initTrace()(err error) {
	serviceConf := commonService.serviceConf
	traceConf := serviceConf.Trace
	if !traceConf.SwitchOn {
		return
	}
	err = trace.Init(serviceConf.ServiceName,traceConf.ReportAddr,traceConf.SampleType,traceConf.SampleRate)
	if err != nil {
		return
	}
	return
}

//服务中间件
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
