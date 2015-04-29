/**
 * Copyright 2015 @ S1N1 Team.
 * name : redis_storage.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package storage

import (
	"github.com/atnet/gof"
	"github.com/garyburd/redigo/redis"
	"fmt"
	"bytes"
	"encoding/gob"
)

type redisStorage struct {
	_pool *redis.Pool
}

func NewRedisStorage(pool *redis.Pool) gof.Storage {
	return &redisStorage{
		_pool: pool,
	}
}

func (this *redisStorage) Get(key string, dst interface{}) error {
	src, err := redis.Values(this._pool.Get().Do("Get", key))
	fmt.Println("--------",src)
	if _, err = redis.Scan(src, &dst); err != nil {
		return err
	}
	return nil
}
func (this *redisStorage) Set(key string, v interface{})error {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(v)
	if err != nil {
		return err
	}
	b := buf.Bytes()
	_,err = this._pool.Get().Do("SET", key,b)
	return err
}
func (this *redisStorage) DSet(key string, v interface{}, seconds int32) {
	this._pool.Get().Do("SETEX", key, v, seconds)
}
