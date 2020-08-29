package main

import (
	"context"
	"fmt"
	"myRPC/registry/base"
	"myRPC/registry/register"
	"time"
)

func main()  {
	registryBase.InitRegistry()
	tecdRegister,err := registryBase.GetRegistryManager().NewRegister("etcd",map[interface{}]interface{}{})
	if err != nil {
		fmt.Println("初始化失败:",err)
		return
	}
	node := &register.Node{NodeId:2,NodeIp:"127.0.0.2",NodePort:"1000",NodeVersion:100,NodeWeight:1,NodeFuncs:[]string{"func1,func2"}}
	service := &register.Service{
		SvrName:"serB",
		SvrType:1001,
		SvrNodes: []*register.Node{
			node,
		},
	}
	err = tecdRegister.Register(service)
	if err != nil {
		fmt.Println("注册失败")
		return
	}
	for {
		temp,err := tecdRegister.GetService(context.TODO(),"serA")
		if err != nil {
			fmt.Println("获取服务A失败")
			return
		}
		fmt.Println("服务A内容:",temp.SvrName,temp.SvrType)
		for _,v := range temp.SvrNodes{
			fmt.Println("服务A节点内容:",v)
		}

		temp,err = tecdRegister.GetService(context.TODO(),"serB")
		if err != nil {
			fmt.Println("获取服务B失败")
			return
		}
		fmt.Println("服务B内容:",temp.SvrName,temp.SvrType)
		for _,v := range temp.SvrNodes{
			fmt.Println("服务B节点内容:",v)
		}

		time.Sleep(10*time.Second)
	}
}
