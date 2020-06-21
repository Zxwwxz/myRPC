package main

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"net/http"
	"time"
)

func main() {
	hystrix.ConfigureCommand("hystrix_rpc", hystrix.CommandConfig{
		Timeout:               1000,
		MaxConcurrentRequests: 100,
		ErrorPercentThreshold: 25,
	})
	for {
		err := hystrix.Do("get_baidu", func() error {
			_, err := http.Get("https://www.baidu.com/")
			if err != nil {
				//fmt.Println("请求失败：",err)
				return err
			}
			return nil
		}, func(err error) error {
			fmt.Println("熔断：", err)
			return err
		})
		if err == nil {
			fmt.Println("请求成功")
		}
		time.Sleep(time.Millisecond * 100)
	}
}
