package limit

import (
	"golang.org/x/time/rate"
)

type TokenLimit struct {
	qps       float64        //令牌桶qps
	allWater  int        //令牌桶最多能装的水大小
	limiter   *rate.Limiter //开源限流器
}

func NewTokenLimit(qps float64, allWater int) *TokenLimit {
	limiter := rate.NewLimiter(rate.Limit(qps), allWater)
	return &TokenLimit{
		qps:        qps,
		allWater:   allWater,
		limiter:    limiter,
	}
}

func (c *TokenLimit) Allow() bool {
	return c.limiter.Allow()
}