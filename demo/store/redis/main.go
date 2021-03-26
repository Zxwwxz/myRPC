package main

import (
    "fmt"
    "github.com/garyburd/redigo/redis"
)

func main() {
    var pool *redis.Pool  //创建redis连接池
    pool = &redis.Pool{     //实例化一个连接池
        //最初的连接数量
        MaxIdle:2,
        //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
        MaxActive:0,
        //连接关闭时间 300秒 （300秒不使用自动关闭）
        IdleTimeout:300,
        //要连接的redis数据库
        Dial: func() (redis.Conn ,error){
            return redis.Dial("tcp","47.92.212.70:6379")
        },
    }
    //连接redis
    c := pool.Get()
    defer c.Close()
    setRedis(c)
    getRedis(c)
}

func setRedis(c redis.Conn){
    //执行redis命令
    _, err := c.Do("Set", "gotest", 200)
    if err != nil {
        fmt.Println("set redis:", err)
        return
    }
}

func getRedis(c redis.Conn){
    //执行redis命令
    r, err := redis.Int(c.Do("Get", "gotest"))
    if err != nil {
        fmt.Println("get redis:", err)
        return
    }
    fmt.Println("get redis result:",r)
}