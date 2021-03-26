package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main()  {
	//创建etcd客户端
	client,err := clientv3.New(clientv3.Config{
		Endpoints:[]string{"47.92.212.70:2379"},
		DialTimeout:time.Second,
	})
	if err != nil {
		fmt.Println("clientv3.New err:",err)
		return
	}
	defer client.Close()
	Watch(client)
}

//设置永久模式
func GetAndPut(client *clientv3.Client)  {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err := client.Put(ctx, "gotest1", "111")
	cancel()
	if err != nil {
		fmt.Println("client.Put err:",err)
		return
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := client.Get(ctx, "gotest1")
	cancel()
	if err != nil {
		fmt.Println("client.Put err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}
}

//设置过期模式
func GrantGetAndPut(client *clientv3.Client)  {
	//申请一个租期，10s过期
	resp ,err := client.Grant(context.TODO(),10)
	if err != nil {
		fmt.Println("client.Grant err:",err)
	}
	_,err = client.Put(context.TODO(),"gotest2","222",clientv3.WithLease(resp.ID))
	if err != nil {
		fmt.Println("client.Put err:",err)
	}
	//续租
	ch ,err := client.KeepAlive(context.TODO(),resp.ID)
	if err != nil {
		fmt.Println("client.KeepAlive err:",err)
	}
	for {
		ka := <- ch
		fmt.Println("续租结果：",ka)
	}
}

//监听模式
func Watch(client *clientv3.Client) {
	rch := client.Watch(context.Background(), "gotest3")
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}
