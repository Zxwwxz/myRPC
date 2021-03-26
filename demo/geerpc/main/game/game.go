package main

import (
    "fmt"
    "geerpc"
    "net"
    "strconv"
)

type GameServer struct {}

type StartGameReq struct{ User int }
type StartGameRsp struct{ StartResult string }

func (g *GameServer) StartGame(req StartGameReq,rsp *StartGameRsp) error {
    fmt.Println("enter StartGame")
    user := req.User
    rsp.StartResult = "user:" + strconv.Itoa(user) + ",start game"
    return nil
}

func main()  {
    err := geerpc.Register(&GameServer{})
    if err != nil {
        fmt.Println("Register err:",err)
    }
    listen, _ := net.Listen("tcp", "127.0.0.1:9999")
    fmt.Println("start accept")
    geerpc.Accept(listen)
}


