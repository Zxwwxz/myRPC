package limiter

import (
	"math"
	"sync"
	"time"
)

const (
	default_bucket_rate = 10
	default_bucket_all_water = 100
)
//漏桶限流
type BucketLimit struct {
	rate       float64 //漏桶中水的漏出速率（每秒）
	curWater   float64 //当前桶里面的水
	allWater   float64 //漏桶最多能装的水大小
	unixNano   int64   //unix时间戳
	lock       sync.RWMutex
}

func NewBucketLimit(params map[interface{}]interface{}) *BucketLimit {
	rate := params["rate"].(int)
	if rate == 0 {
		rate = default_bucket_rate
	}
	allWater := params["all_water"].(int)
	if rate == 0 {
		allWater = default_bucket_all_water
	}
	bucketLimit := &BucketLimit{
		rate:   	float64(rate),
		curWater:   0,
		allWater:   float64(allWater),
		unixNano:   0,
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

//是否通过
func (b *BucketLimit) Allow() bool {
	b.lock.Lock()
	defer b.lock.Unlock()
	//先漏水
	b.reflesh()
	if b.curWater < b.allWater {
		//加水
		b.curWater = b.curWater + 1
		return true
	}
	return false
}