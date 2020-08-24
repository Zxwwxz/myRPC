
package main

import (
	"fmt"
	pbGame "myRPC/myService/game/generate"
	"myRPC/myService/game/router"
	"myRPC/service"
)

func main()  {
	err := service.InitService()
	if err != nil {
		fmt.Println("init server err:",err)
	}
	pbGame.RegisterGameServiceServer(service.GetGrpcService(),&router.Router{})
	fmt.Println("init server success")	
	service.Run()
}

