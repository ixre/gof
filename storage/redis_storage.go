/**
 * Copyright 2015 @ z3q.net.
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
	"github.com/garyburd/redigo/redis"
	"github.com/jsix/gof"
	"reflect"
	"strings"
	"sync"
)

var DriveRedisStorage string = "redis-storage"

type redisStorage struct {
	_pool *redis.Pool
	_buf  *bytes.Buffer
	_conn redis.Conn
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

func isBaseOfStruct(v interface{}) bool {
	valueType := reflect.TypeOf(v)
	kind := valueType.Kind()
	if kind == reflect.Ptr {
		kind = valueType.Elem().Kind()
	}
	return kind == reflect.Struct || kind == reflect.Map || kind == reflect.Array
}

func (this *redisStorage) getRedisBytes(key string) ([]byte, error) {
	src, err := redis.Bytes(this.getConn().Do("GET", key))
	return src, err
}

// return storage driver
func (this *redisStorage) Driver() interface{} {
	return this._pool
}

func (this *redisStorage) DriverName() string {
	return DriveRedisStorage
}

func (this *redisStorage) Exists(key string) (exists bool) {
	i, err := redis.Int(this.getConn().Do("EXISTS", key))
	return err != nil && i == 1
}

func (this *redisStorage) Get(key string, dst interface{}) error {
	if isBaseOfStruct(dst) {
		src, err := this.getRedisBytes(key)
		if err == nil {
			err = this.decodeBytes(src, dst)
		}
		return err
	}
	return errors.New("dst must be struct")
}

func (this *redisStorage) getConn() redis.Conn {
	if this._conn == nil {
		this._conn = this._pool.Get()
	}
	return this._conn
}
func (this *redisStorage) GetBool(key string) (bool, error) {
	src, err := redis.Bool(this.getConn().Do("GET", key))
	return src, err
}

func (this *redisStorage) GetInt(key string) (int, error) {
	src, err := redis.Int(this.getConn().Do("GET", key))
	return src, err
}

func (this *redisStorage) GetInt64(key string) (int64, error) {
	src, err := redis.Int64(this.getConn().Do("GET", key))
	return src, err
}

func (this *redisStorage) GetString(key string) (string, error) {
	d, err := this.getRedisBytes(key)
	if err != nil {
		return "", err
	}
	return string(d), err
}

func (this *redisStorage) GetFloat64(key string) (float64, error) {
	src, err := redis.Float64(this.getConn().Do("GET", key))
	return src, err
}

//Get raw value
func (this *redisStorage) GetRaw(key string) (interface{}, error) {
	replay, err := this.getConn().Do("GET", key)
	return replay, err
}

func (this *redisStorage) Set(key string, v interface{}) error {
	var err error
	var redisValue interface{} = v
	if isBaseOfStruct(v) {
		redisValue, err = this.getByte(v)
	}
	_, err = this.getConn().Do("SET", key, redisValue)
	return err
}
func (this *redisStorage) Del(key string) {
	this.getConn().Do("DEL", key)
}

func (this *redisStorage) SetExpire(key string, v interface{}, seconds int64) error {
	var err error
	var redisValue interface{} = v
	if isBaseOfStruct(v) {
		redisValue, err = this.getByte(v)
	}
	_, err = this.getConn().Do("SETEX", key, seconds, redisValue)
	return err
}
