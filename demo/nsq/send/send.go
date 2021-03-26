//组件：
//nsqlookupd：管理nsqd,nsqadmin拓扑信息
//nsqlookupd
//tcp监听地址4160，http监听地址4161
//nsqd：接收消息，排队消息，投递消息，生产者消费者连接对象
//nsqd --lookupd-tcp-address=127.0.0.1:4160 -tcp-address="0.0.0.0:4150 -broadcast-address=47.92.212.70"
//tcp监听地址4150，http监听地址4151
//nsqadmin：web
//nsqadmin --lookupd-http-address=127.0.0.1:4161
//http监听地址4171

//流程：
//生产者生成msg-》nsqd的topic中
//消费者向nsqlookupd查询带有topic的nsqd，nsqd新建channel与topic建立连接
//nsqd推送msg给订阅topic给所有channel，其实就是推给了消费者

//Nsq发送测试
package main

import (
    "github.com/nsqio/go-nsq"
    "log"
)


func main() {
    config := nsq.NewConfig()
    //创建生产者
    w, _ := nsq.NewProducer("47.92.212.70:4150", config)
    //向nsqd的topic发送消息
    err := w.Publish("test_topic3", []byte("test_msg1"))
    if err != nil {
        log.Panic("Could not connect")
    }
    w.Stop()
}