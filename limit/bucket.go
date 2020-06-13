package limit

import (
	"math"
	"time"
)

type BucketLimit struct {
	rate       float64 //漏桶中水的漏出速率
	allWater   float64 //漏桶最多能装的水大小
	unixNano   int64   //unix时间戳
	curWater   float64 //当前桶里面的水
}

func NewBucketLimit(rate float64, allWater int64) *BucketLimit {
	return &BucketLimit{
		allWater:   float64(allWater),
		rate:       rate,
		unixNano:   time.Now().UnixNano(),
		curWater:   0,
	}
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