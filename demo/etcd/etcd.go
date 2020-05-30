package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main()  {
	client,err := clientv3.New(clientv3.Config{
		Endpoints:[]string{"47.92.212.70:2379"},
		DialTimeout:time.Second,
	})
	if err != nil {
		fmt.Println("err1:",err)
	}
	defer client.Close()
	resp ,err := client.Grant(context.TODO(),10)
	if err != nil {
		fmt.Println("err2:",err)
	}
	_,err = client.Put(context.TODO(),"aaaa","1111",clientv3.WithLease(resp.ID))
	if err != nil {
		fmt.Println("err3:",err)
	}

	ch ,err := client.KeepAlive(context.TODO(),resp.ID)
	if err != nil {
		fmt.Println("err3:",err)
	}
	for {
		ka := <- ch
		fmt.Println("通道：",ka)
	}
}
