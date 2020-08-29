
package main

import (
	"context"
	"fmt"
	"myRPC/log/base"
	"myRPC/meta"
	"myRPC/myService/game/client"
	"myRPC/myService/game/generate"
	"myRPC/myService/login/generate"
	"myRPC/myService/login/router"
	"myRPC/service"
)

func main()  {
	err := service.InitService()
	if err != nil {
		logBase.Fatal("init server err:%v",err)
		fmt.Printf("init server err:%v",err)
	}
	pbLogin.RegisterLoginRpcServer(service.GetGrpcService(),&router.Router{})
	fmt.Printf("init server run")
	logBase.Debug("init server run")
	go func() {
		send()
	}()
	service.Run()
}

func send()  {
	resp ,err := gameClient.NewClientCall().StartGame(context.Background(),&pbGame.StartGameRequest{StartName:"aaa"},[]meta.ClientMetaOption{
		meta.SetMaxReconnectNum(3),
		meta.SetCallerType(meta.Caller_type_balance),
	})
	if err != nil {
		fmt.Println("rpc call err:",err)
	}
	fmt.Println("resp:",resp)
}

