
package main

import (
	"myRPC/log/base"
    "myRPC/service"
	"myRPC/myService/game/generate"
	"myRPC/myService/game/router"
	
)

func main()  {
	err := service.InitService()
	if err != nil {
		logBase.Fatal("init server err:%v",err)
	}
	pbGame.RegisterGameServiceServer(service.GetGrpcService(),&router.Router{})
	logBase.Debug("init server run")
	service.Run()
}

