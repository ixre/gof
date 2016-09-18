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
	"github.com/jsix/gof/log"
	"reflect"
	"strings"
)

type IRedisStorage interface {
	// get keys start with prefix
	Keys(prefix string) ([]string, error)
	// delete keys contain prefix
	PrefixDel(prefix string) (int, error)
}

var DriveRedisStorage string = "redis-storage"
var _ Interface = new(redisStorage)
var _ IRedisStorage = new(redisStorage)

func NewRedisStorage(pool *redis.Pool) Interface {
	return &redisStorage{
		pool: pool,
	}
}

type redisStorage struct {
	pool *redis.Pool
}

func (r *redisStorage) encodeBytes(v interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(v)
	r.checkTypeErr(err)
	return buf.Bytes(), err
}

func (r *redisStorage) decodeBytes(b []byte, dst interface{}) error {
	buf := bytes.NewBuffer(b)
	err := gob.NewDecoder(buf).Decode(dst)
	log.Println("--- redis ", err)
	r.checkTypeErr(err)
	return err
}

func (r *redisStorage) checkTypeErr(err error) {
	if err != nil && strings.Index(err.Error(), "type not registered") != -1 {
		panic(err)
	}
}

func checkInputValueType(v interface{}) bool {
	kind := reflect.TypeOf(v).Kind()
	return kind == reflect.Ptr || kind == reflect.Struct ||
		kind == reflect.Map || kind == reflect.Array
}

func checkOutputValueType(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Ptr
	//vType := reflect.TypeOf(v)
	//kind := vType.Kind()
	//if kind == reflect.Ptr {
	//	kind = vType.Elem().Kind()
	//	if kind == reflect.Ptr {
	//		panic(errors.New("dst ptr is a ptr."))
	//	}
	//}
	//return kind == reflect.Ptr
}

func (r *redisStorage) Driver() interface{} {
	return r.pool
}

func (r *redisStorage) DriverName() string {
	return DriveRedisStorage
}

func (r *redisStorage) getBytes(key string) ([]byte, error) {
	conn := r.pool.Get()
	defer conn.Close()
	src, err := redis.Bytes(conn.Do("GET", key))
	return src, err
}

func (r *redisStorage) Get(key string, dst interface{}) error {
	if checkOutputValueType(dst) {
		src, err := r.getBytes(key)
		if err == nil {
			err = r.decodeBytes(src, dst)
		}
		return err
	}
	return errors.New("dst must be struct")
}

func (r *redisStorage) GetBool(key string) (bool, error) {
	conn := r.pool.Get()
	defer conn.Close()
	src, err := redis.Bool(conn.Do("GET", key))
	return src, err
}

func (r *redisStorage) GetInt(key string) (int, error) {
	conn := r.pool.Get()
	defer conn.Close()
	src, err := redis.Int(conn.Do("GET", key))
	return src, err
}

func (r *redisStorage) GetInt64(key string) (int64, error) {
	conn := r.pool.Get()
	defer conn.Close()
	src, err := redis.Int64(conn.Do("GET", key))
	return src, err
}

func (r *redisStorage) GetString(key string) (string, error) {
	d, err := r.getBytes(key)
	if err != nil {
		return "", err
	}
	return string(d), err
}

func (r *redisStorage) GetFloat64(key string) (float64, error) {
	conn := r.pool.Get()
	defer conn.Close()
	src, err := redis.Float64(conn.Do("GET", key))
	return src, err
}

func (r *redisStorage) GetRaw(key string) (interface{}, error) {
	conn := r.pool.Get()
	defer conn.Close()
	replay, err := conn.Do("GET", key)
	return replay, err
}

func (r *redisStorage) Exists(key string) bool {
	conn := r.pool.Get()
	defer conn.Close()
	i, err := redis.Int(conn.Do("EXISTS", key))
	return err != nil && i == 1
}

func (r *redisStorage) Set(key string, v interface{}) error {
	var err error
	var redisValue interface{} = v
	if checkInputValueType(v) {
		redisValue, err = r.encodeBytes(v)
	}
	conn := r.pool.Get()
	defer conn.Close()
	_, err = conn.Do("SET", key, redisValue)
	return err
}

func (r *redisStorage) Del(key string) {
	conn := r.pool.Get()
	defer conn.Close()
	conn.Do("DEL", key)
}

func (r *redisStorage) SetExpire(key string, v interface{}, seconds int64) error {
	var err error
	var redisValue interface{} = v
	if checkInputValueType(v) {
		redisValue, err = r.encodeBytes(v)
	}
	conn := r.pool.Get()
	defer conn.Close()
	_, err = conn.Do("SETEX", key, seconds, redisValue)
	return err
}

// get keys start with prefix
func (r *redisStorage) Keys(prefix string) ([]string, error) {
	conn := r.pool.Get()
	defer conn.Close()
	return redis.Strings(conn.Do("KEYS", prefix))
}

// delete keys contain prefix
func (r *redisStorage) PrefixDel(prefix string) (int, error) {
	keys, err := r.Keys(prefix)
	if err != nil {
		return 0, err
	}
	conn := r.pool.Get()
	defer conn.Close()
	for _, key := range keys {
		conn.Do("DEL", key)
	}
	return len(keys), nil
}
