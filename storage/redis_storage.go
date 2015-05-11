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
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/atnet/gof"
	"github.com/garyburd/redigo/redis"
	"strings"
	"sync"
)

var DriveRedisStorage string = "redis-storage"

type redisStorage struct {
	_pool *redis.Pool
	_buf  *bytes.Buffer
	sync.Mutex
}

func NewRedisStorage(pool *redis.Pool) gof.Storage {
	return &redisStorage{
		_pool: pool,
		_buf:  new(bytes.Buffer),
	}
}

func (this *redisStorage) getByte(v interface{}) ([]byte, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	enc := gob.NewEncoder(this._buf)
	err := enc.Encode(v)
	if err == nil {
		b := this._buf.Bytes()
		this._buf.Reset()
		return b, nil
	}
	if strings.Index(err.Error(), "type not registered") != -1 {
		panic(err)
	}
	return nil, err
}

func (this *redisStorage) decodeBytes(b []byte, dst interface{}) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	this._buf.Write(b)
	dec := gob.NewDecoder(this._buf)
	err := dec.Decode(dst)
	this._buf.Reset()
	return err
}

func (this *redisStorage) Driver() string {
	return DriveRedisStorage
}

func (this *redisStorage) Get(key string, dst interface{}) error {
	conn := this._pool.Get()
	src, err := redis.Bytes(conn.Do("GET", key))
	conn.Close()
	if err == nil{
		err = this.decodeBytes(src, dst)
	}
	return err
}

//Get raw value
func (this *redisStorage) GetRaw(key string) interface{} {
	panic(errors.New("HashStorage not support method \"GetRaw\""))
}

func (this *redisStorage) Set(key string, v interface{}) error {
	if v != nil {
		b, err := this.getByte(v)

		if err == nil {
			conn := this._pool.Get()
			_, err = conn.Do("SET", key, b)
			conn.Close()
		}
		return err
	}
	return errors.New("value can't be nil.")
}

func (this *redisStorage) Del(key string) {
	conn := this._pool.Get()
	conn.Do("DEL", key)
}

func (this *redisStorage) SetExpire(key string, v interface{}, seconds int64) error {
	b, err := this.getByte(v)
	if err == nil {
		conn := this._pool.Get()
		_, err = this._pool.Get().Do("SETEX", key, seconds, b)
		conn.Close()
	}
	return err
}
