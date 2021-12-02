package concurrent

import (
	"testing"
	"time"
)

/**
 * Copyright 2009-2019 @ 56x.net
 * name : token_bucket_test.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-10-01 17:29
 * description :
 * history :
 */

// 每秒最多请求10个令牌
func TestTokenBucket_Acquire(t *testing.T) {
	bucket := NewTokenBucket(100, 10)
	for {
		for i := 0; i < 100; i++ {
			b := bucket.Acquire(1)
			println("--- Req:", i, " => ", b)
		}
		time.Sleep(time.Second)
	}
}
