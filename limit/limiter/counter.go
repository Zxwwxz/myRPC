package limiter

import (
	"sync/atomic"
	"time"
)

const (
	default_counter_limit_num = 10
	default_counter_interval_nano = 10
)

type CounterLimit struct {
	counterNum      int64 //计数器
	limitNum        int64 //指定时间窗口内允许的最大请求数
	intervalNano    int64 //指定时间窗口
	lastNano        int64 //unix时间戳,单位为纳秒
}

func NewCounterLimit(params map[interface{}]interface{}) *CounterLimit {
	counterLimit := &CounterLimit{
		counterNum:   	0,
		limitNum:   	default_counter_limit_num,
		intervalNano:   default_counter_interval_nano,
		lastNano:       0,
	}
	limitNum := params["limit_num"].(int64)
	if limitNum != 0 {
		counterLimit.limitNum = limitNum
	}
	intervalNano := params["interval_nano"].(int64)
	if intervalNano != 0 {
		counterLimit.intervalNano = intervalNano
	}
	return counterLimit
}

func (c *CounterLimit) Allow() bool {
	now := time.Now().UnixNano()
	if now-c.lastNano > c.intervalNano {
		atomic.StoreInt64(&c.counterNum, 0)
		atomic.StoreInt64(&c.lastNano, now)
		return true
	}
	atomic.AddInt64(&c.counterNum, 1)
	return c.counterNum < c.limitNum
}
