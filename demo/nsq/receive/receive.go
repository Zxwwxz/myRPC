package main

import (
    "fmt"
    "github.com/nsqio/go-nsq"
    "log"
    "sync"
)

func main() {
    wg := &sync.WaitGroup{}
    wg.Add(1)

    config := nsq.NewConfig()
    //新建消费者，指定监听的topic，指定channel名
    q, _ := nsq.NewConsumer("test_topic3", "test_channel", config)
    //注册回调函数，topic有收到消息，会通过channel发送到这里
    q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
        fmt.Println(string(message.Body))
        wg.Done()
        return nil
    }))
    //向nsqlookupd查找topic的nsqd，并连接nsqd
    err := q.ConnectToNSQLookupd("47.92.212.70:4161")
    if err != nil {
        log.Panic("Could not connect")
    }
    wg.Wait()
}