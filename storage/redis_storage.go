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
	"bytes"
	"encoding/gob"
    "sync"
)

type redisStorage struct {
	_pool *redis.Pool
	_buf *bytes.Buffer
    sync.Mutex
}

func NewRedisStorage(pool *redis.Pool) gof.Storage {
	return &redisStorage{
		_pool: pool,
        _buf : new(bytes.Buffer),
	}
}

func (this *redisStorage) Get(key string, dst interface{}) error {
	src, err := redis.Bytes(this._pool.Get().Do("Get", key))
	if err == nil{
        err = this.decodeBytes(src,dst)
	}
	return err
}

func (this *redisStorage) getByte(v interface{})([]byte,error){
    this.Mutex.Lock()
    defer this.Mutex.Unlock()
	enc := gob.NewEncoder(this._buf)
	err := enc.Encode(v)
	if err == nil {
		return nil,err
	}
    b := this._buf.Bytes()
    this._buf.Reset()
	return b,err
}

func (this *redisStorage) decodeBytes(b []byte,dst interface{})error {
    this.Mutex.Lock()
    defer this.Mutex.Unlock()
    dec := gob.NewDecoder(this._buf)
    err := dec.Decode(dst)
    this._buf.Reset()
    return err
}

func (this *redisStorage) Set(key string, v interface{})error {
	b,err := this.getByte(v)
	if err == nil {
		_, err = this._pool.Get().Do("SET", key, b)
	}
	return err
}
func (this *redisStorage) SetExpire(key string, v interface{}, seconds int32)error {
	b,err := this.getByte(v)
	if err == nil {
		_, err = this._pool.Get().Do("SETEX", key, b,seconds)
	}
	return err
}
