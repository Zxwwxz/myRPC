package service

import (
	"fmt"
	"google.golang.org/grpc"
	"myRPC/config"
	"myRPC/limit"
	mwBase "myRPC/middleware/base"
	mwLimit "myRPC/middleware/limit"
	mwPrometheus "myRPC/middleware/prometheus"
	"net"
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

func BuildServerMiddleware(handle mwBase.MiddleWareFunc,frontMiddles,backMiddles []mwBase.MiddleWare) mwBase.MiddleWareFunc {
	var middles []mwBase.MiddleWare
	serviceConf := config.GetConf()
	middles = append(middles,frontMiddles...)
	if serviceConf.Prometheus.SwitchOn {
		middles = append(middles, mwPrometheus.PrometheusMiddleware())
	}
	if serviceConf.Limit.SwitchOn {
		middles = append(middles, mwLimit.LimitMiddleware(commonService.Limiter))
	}
	middles = append(middles,backMiddles...)
	m := mwBase.Chain(middles...)
	return m(handle)
}
