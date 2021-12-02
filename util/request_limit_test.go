package util

import (
	"github.com/ixre/gof/storage"
	"testing"
	"time"
)

/**
 * Copyright 2009-2019 @ 56x.net
 * name : request_limit_test.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-10-02 11:02
 * description :
 * history :
 */

func TestRequestLimit_Acquire(t *testing.T) {
	rl := NewRequestLimit(storage.NewHashStorage(), 100, 10, 30)
	ip := "172.17.0.1"
	for {
		for i := 0; i < 100; i++ {
			if rl.IsLock(ip) {
				println("ip locked,please try later")
				continue
			}
			b := rl.Acquire(ip, 1)
			println("--- Req:", i, " => ", b)
		}
		time.Sleep(time.Second)
	}
}
