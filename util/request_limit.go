package util

import (
	"fmt"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/util/concurrent"
	"sync"
)

/**
 * Copyright 2009-2019 @ to2.net
 * name : request_limit.go.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-10-01 16:48
 * description :
 * history :
 */

type RequestLimit struct {
	buckets map[string]*concurrent.TokenBucket
	sync.RWMutex
	capacity   int     // 桶的容量
	rate       float64 // 令牌放入速度
	lockSecond int64
	store      storage.Interface
}

// 创建请求限制, store存储数据,lockSecond锁定时间,单位:秒,capacity: 最大容量,rate: 令牌放入速度
func NewRequestLimit(store storage.Interface, capacity int, rate float64, lockSecond int) *RequestLimit {
	return &RequestLimit{
		buckets:    make(map[string]*concurrent.TokenBucket, 0),
		store:      store,
		RWMutex:      sync.RWMutex{},
		capacity:   capacity,
		rate:       rate,
		lockSecond: int64(lockSecond),
	}
}

// 是否锁定
func (i *RequestLimit) IsLock(addr string) bool {
	k := fmt.Sprintf("sys:req-limit:%s", addr)
	v, err := i.store.GetInt(k)
	return err == nil && v > 0
}

// 锁定地址
func (i *RequestLimit) lockAddr(addr string) {
	k := fmt.Sprintf("sys:req-limit:%s", addr)
	_ = i.store.SetExpire(k, 1, i.lockSecond)
}

// Acquire 获取令牌
func (i *RequestLimit) Acquire(addr string, n int) bool {
	i.RWMutex.RLock()
	v, ok := i.buckets[addr]
	i.RWMutex.RUnlock()
	if !ok {
		v = concurrent.NewTokenBucket(i.capacity, i.rate)
		i.RWMutex.Lock()
		i.buckets[addr] = v
		i.RWMutex.Unlock()
	}
	b := v.Acquire(n)
	if !b {
		i.lockAddr(addr)
	}
	return b
}
