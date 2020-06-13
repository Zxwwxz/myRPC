package limit

import (
	"sync/atomic"
	"time"
)

type CounterLimit struct {
	counterNum      int64 //计数器
	limitNum        int64 //指定时间窗口内允许的最大请求数
	intervalNano    int64 //指定的时间窗口
	lastNano        int64 //unix时间戳,单位为纳秒
}
func NewCounterLimit(interval time.Duration, limit int64) *CounterLimit {
	return &CounterLimit{
		counterNum:      0,
		limitNum:        limit,
		intervalNano:    int64(interval),
		lastNano:        time.Now().UnixNano(),
	}
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
