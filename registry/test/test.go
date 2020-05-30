package main

import (
	"context"
	"fmt"
	"github.com/ibinarytree/koala/registry/base"
	"github.com/ibinarytree/koala/registry/etcd"
	"time"
)

func main()  {
	tecdPlugin,err := base.PluginManager.InitPlugin(context.TODO(),etcd.EtcdPluginName,
		base.SetRegisterAddrs([]string{"47.92.212.70:2379"}),
		base.SetRegisterPath("/myRPC"),
		base.SetRegisterTimeOut(2 * time.Second),
		base.SetHeartTimeOut(5))
	if err != nil {
		fmt.Println("初始化失败:",err)
		return
	}
	node := &base.Node{NodeId:2,NodeIp:"127.0.0.2",NodePort:"1000",NodeVersion:100,NodeWeight:1,NodeFuncs:[]string{"func1,func2"}}
	service := &base.Service{
		SvrName:"serB",
		SvrType:1001,
		SvrNodes: map[int]*base.Node{
			2:node,
		},
	}
	err = tecdPlugin.Register(context.TODO(),service)
	if err != nil {
		fmt.Println("注册失败")
		return
	}
	for {
		temp,err := tecdPlugin.GetService(context.TODO(),"serA")
		if err != nil {
			fmt.Println("获取服务A失败")
			return
		}
		fmt.Println("服务A内容:",temp.SvrName,temp.SvrType)
		for _,v := range temp.SvrNodes{
			fmt.Println("服务A节点内容:",v)
		}

		temp,err = tecdPlugin.GetService(context.TODO(),"serB")
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
