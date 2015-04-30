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
    "strings"
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

func (this *redisStorage) getByte(v interface{})([]byte,error) {
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

func (this *redisStorage) decodeBytes(b []byte,dst interface{})error {
    this.Mutex.Lock()
    defer this.Mutex.Unlock()
    this._buf.Write(b)
    dec := gob.NewDecoder(this._buf)
    err := dec.Decode(dst)
    this._buf.Reset()
    return err
}


func (this *redisStorage) Get(key string, dst interface{}) error {
    conn := this._pool.Get()
    src, err := redis.Bytes(conn.Do("GET", key))
    conn.Close()
    if err == nil{
        err = this.decodeBytes(src,dst)
    }
    return err
}


func (this *redisStorage) Set(key string, v interface{})error {
	b,err := this.getByte(v)
	if err == nil {
        conn := this._pool.Get()
		_, err = conn.Do("SET", key, b)
        conn.Close()
	}
	return err
}

func (this *redisStorage) Del(key string){
    conn := this._pool.Get()
    conn.Do("DEL",key)
}

func (this *redisStorage) SetExpire(key string, v interface{}, seconds int32)error {
	b,err := this.getByte(v)
	if err == nil {
        conn := this._pool.Get()
        _, err = this._pool.Get().Do("SETEX", key, b,seconds)
        conn.Close()
	}
	return err
}
