package limiter

import (
	"math"
	"time"
)

const (
	default_bucket_rate = 10
	default_bucket_cur_water = 10
	default_bucket_all_water = 10
)

type BucketLimit struct {
	rate       float64 //漏桶中水的漏出速率
	curWater   float64 //当前桶里面的水
	allWater   float64 //漏桶最多能装的水大小
	unixNano   int64   //unix时间戳
}

func NewBucketLimit(params map[interface{}]interface{}) *BucketLimit {
	bucketLimit := &BucketLimit{
		rate:   	default_bucket_rate,
		curWater:   default_bucket_cur_water,
		allWater:   default_bucket_all_water,
		unixNano:   0,
	}
	rate := params["rate"].(float64)
	if rate != 0 {
		bucketLimit.rate = rate
	}
	curWater := params["cur_water"].(float64)
	if rate != 0 {
		bucketLimit.curWater = curWater
	}
	allWater := params["all_water"].(float64)
	if rate != 0 {
		bucketLimit.allWater = allWater
	}
	return bucketLimit
}

func (b *BucketLimit) reflesh() {
	now := time.Now().UnixNano()
	diffSec := float64(now-b.unixNano) / 1000 / 1000 / 1000
	b.curWater = math.Max(0, b.curWater-diffSec*b.rate)
	b.unixNano = now
	return
}

func (b *BucketLimit) Allow() bool {
	b.reflesh()
	if b.curWater < b.allWater {
		b.curWater = b.curWater + 1
		return true
	}
	return false
}