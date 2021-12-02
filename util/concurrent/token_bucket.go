package concurrent

import (
	"time"
)

/**
 * Copyright 2009-2019 @ 56x.net
 * name : token_bucket.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-10-01 16:48
 * description :
 * history :
 */

type TokenBucket struct {
	timestamp int64   // 时间
	capacity  int     // 桶的容量
	rate      float64 // 令牌放入速度
	tokens    int     // 当前令牌数量
}

// capacity 令牌桶的容量; rate 放入令牌桶的速度/每秒放入令牌的数量
func NewTokenBucket(capacity int, rate float64) *TokenBucket {
	return &TokenBucket{
		timestamp: 0,
		capacity:  capacity,
		rate:      rate,
		tokens:    capacity,
	}
}

// 获取n个令牌
func (t *TokenBucket) Acquire(n int) bool {
	now := time.Now().Unix()
	// 计算当前的令牌数
	t.tokens = t.tokens + int(float64(now-t.timestamp)*t.rate)
	// 保存上次请求的时间戳
	t.timestamp = now
	// 如果令牌超出令牌桶的容量
	if t.tokens > t.capacity {
		t.tokens = t.capacity
	}
	// 获取令牌,并减去令牌数
	if t.tokens > n {
		t.tokens -= n
		return true
	}
	return false
}
