package main

import (
	"fmt"
	"myRPC/limit/base"
	"time"
)

func main() {
	limitBase.InitLimit()
	limiter,_ := limitBase.GetLimitMgr().NewLimiter("counter", map[interface{}]interface{}{})
	m := make(map[int]bool)
	for i := 0; i < 1000; i++ {
		allow := limiter.Allow()
		if allow {
			m[i] = true
		} else {
			m[i] = false
		}
		time.Sleep(time.Millisecond*3)
	}
	for i := 0; i < 1000; i++ {
		fmt.Printf("i=%d allow=%v\n", i, m[i])
	}
}
