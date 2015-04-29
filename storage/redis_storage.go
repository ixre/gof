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
	_buf bytes.Buffer
}

func NewRedisStorage(pool *redis.Pool) gof.Storage {
	return &redisStorage{
		_pool: pool,
	}
}

func (this *redisStorage) Get(key string, dst interface{}) error {
	src, err := redis.Bytes(this._pool.Get().Do("Get", key))
	if err == nil{
		buf :=
		dec := gob.Decoder(buf)
		buf.Reset()
	}
	return err
}

func getByte(v interface{})([]byte,error){
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(v)
	if err == nil {
		return nil,err
	}
	return buf.Bytes(),err
}
func (this *redisStorage) Set(key string, v interface{})error {
	b,err := getByte(v)
	if err == nil {
		_, err = this._pool.Get().Do("SET", key, b)
	}
	return err
}
func (this *redisStorage) SetExpire(key string, v interface{}, seconds int32)error {
	b,err := getByte(v)
	if err == nil {
		_, err = this._pool.Get().Do("SETEX", key, b,seconds)
	}
	return err
}
