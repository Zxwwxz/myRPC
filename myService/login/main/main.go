
package main

import (
	"fmt"
	pbLogin "myRPC/myService/login/generate"
	"myRPC/myService/login/router"
	"myRPC/service"
)

func main()  {
	err := service.Init()
	if err != nil {
		fmt.Println("init server err:",err)
	}
	pbLogin.RegisterLoginRpcServer(service.GetGrpcService(),&router.Router{})
	fmt.Println("init server success")	
	service.Run()
}

