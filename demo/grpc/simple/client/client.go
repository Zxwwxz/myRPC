package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"myRPC/myService/game/generate"
)

//func main()  {
//	req := &pbGuide.Point{
//		Latitude:100,
//		Longitude:10000,
//	}
//	conn, err := grpc.Dial("127.0.0.1:8889", grpc.WithInsecure())
//	if err != nil {
//		fmt.Println("Dial:", err)
//		return
//	}
//	// 函数结束时关闭连接
//	defer conn.Close()
//	client := pbGuide.NewRouteGuideClient(conn)
//	resp,err := client.GetFeature(context.TODO(),req)
//	if err != nil {
//		fmt.Println("GetFeature:", err)
//		return
//	}
//	fmt.Printf("resp:%v",resp)
//}


func main()  {
	req := &pbGame.StartGameRequest{
		StartName:"",
	}
	conn, err := grpc.Dial("47.92.212.70:8888", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Dial:", err)
		return
	}
	// 函数结束时关闭连接
	defer conn.Close()
	client := pbGame.NewGameServiceClient(conn)
	resp,err := client.StartGame(context.TODO(),req)
	if err != nil {
		fmt.Println("StartGame err:", err)
		return
	}
	fmt.Printf("resp:%v",resp)
}