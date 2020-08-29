package limiter

import (
	"sync"
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
	lock            sync.RWMutex
}

func NewCounterLimit(params map[interface{}]interface{}) *CounterLimit {

	limitNum := params["limit_num"].(int64)
	if limitNum == 0 {
		limitNum = default_counter_limit_num
	}
	intervalNano := params["interval_nano"].(int64)
	if intervalNano == 0 {
		intervalNano = default_counter_interval_nano
	}
	counterLimit := &CounterLimit{
		counterNum:   	0,
		limitNum:   	limitNum,
		intervalNano:   intervalNano,
		lastNano:       0,
	}
	return counterLimit
}

func (c *CounterLimit) Allow() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	now := time.Now().UnixNano()
	if now-c.lastNano > c.intervalNano {
		c.counterNum = 0
		c.lastNano = now
		return true
	}
	c.counterNum = c.counterNum + 1
	return c.counterNum < c.limitNum
}
