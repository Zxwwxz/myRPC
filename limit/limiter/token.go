package limiter

import (
	"golang.org/x/time/rate"
)

const (
	default_token_qps = 10
	default_token_all_water = 10
)
//令牌桶限流
type TokenLimit struct {
	qps       float64        //令牌桶qps
	allWater  int        	 //令牌桶最多能装的水大小
	limiter   *rate.Limiter  //开源限流器
}

func NewTokenLimit(params map[interface{}]interface{}) *TokenLimit {
	qps := params["qps"].(int)
	if qps == 0 {
		qps = default_token_qps
	}
	allWater := params["all_water"].(int)
	if allWater == 0 {
		allWater = default_token_all_water
	}
	limiter := rate.NewLimiter(rate.Limit(qps), allWater)
	tokenLimiter := &TokenLimit{
		qps:		float64(qps),
		allWater:	allWater,
		limiter:limiter,
	}
	return tokenLimiter
}

func (c *TokenLimit) Allow() bool {
	return c.limiter.Allow()
}