package main

import (
    "context"
    "fmt"
    "geerpc"
)

type StartGameReq struct{ User int }
type StartGameRsp struct{ StartResult string }

func main()  {
    client,err := geerpc.XDial("tcp@127.0.0.1:9999",geerpc.DefaultOption)
    if err != nil {
        fmt.Println("XDial err:",err)
        return
    }
    req := &StartGameReq{User:111}
    rsp := &StartGameRsp{}
    err = client.Call(context.TODO(),"GameServer.StartGame",req,rsp)
    if err != nil {
        fmt.Println("Call err:",err)
        return
    }
    fmt.Println("rsp:",rsp)
}
