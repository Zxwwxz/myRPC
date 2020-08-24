package limiter

import (
	"golang.org/x/time/rate"
)

const (
	default_token_qps = 10
	default_token_all_water = 10
)

type TokenLimit struct {
	qps       float64        //令牌桶qps
	allWater  int        	 //令牌桶最多能装的水大小
	limiter   *rate.Limiter  //开源限流器
}

func NewTokenLimit(params map[interface{}]interface{}) *TokenLimit {
	tokenLimiter := &TokenLimit{
		qps:		default_token_qps,
		allWater:	default_token_all_water,
	}
	qps := params["qps"].(float64)
	if qps != 0 {
		tokenLimiter.qps = qps
	}
	allWater := params["all_water"].(int)
	if allWater != 0 {
		tokenLimiter.allWater = allWater
	}
	limiter := rate.NewLimiter(rate.Limit(tokenLimiter.qps), tokenLimiter.allWater)
	tokenLimiter.limiter = limiter
}

func (c *TokenLimit) Allow() bool {
	return c.limiter.Allow()
}