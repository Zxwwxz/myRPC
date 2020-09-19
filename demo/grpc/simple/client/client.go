package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"myRPC/demo/grpc/simple/pbGuide"
)

func main()  {
	req := &pbGuide.Point{
		Latitude:100,
		Longitude:10000,
	}
	conn, err := grpc.Dial("127.0.0.1:8889", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Dial:", err)
		return
	}
	// 函数结束时关闭连接
	defer conn.Close()
	client := pbGuide.NewRouteGuideClient(conn)
	resp,err := client.GetFeature(context.TODO(),req)
	if err != nil {
		fmt.Println("GetFeature:", err)
		return
	}
	fmt.Printf("resp:%v",resp)
}
