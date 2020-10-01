/**
 * Copyright 2015 @ to2.net.
 * name : redis_storage.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"reflect"
	"time"
)

// Create a pool of Redis client
func NewRedisPool(host string, port int, db int, auth string,
	maxIdle int, idleTimeout int) *redis.Pool {

	if port <= 0 {
		port = 6379
	}
	if maxIdle <= 0 {
		maxIdle = 10000
	}
	if idleTimeout <= 0 {
		idleTimeout = 20000
	}

	return &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			var c redis.Conn
			var err error
			for {
				c, err = redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
				if err == nil {
					break
				}
				log.Printf("[ Redis] - redis(%s:%d) dial failed - %s , Redial after 5 seconds\n",
					host, port, err.Error())
				time.Sleep(time.Second * 5)
			}

			if len(auth) != 0 {
				if _, err = c.Do("AUTH", auth); err != nil {
					c.Close()
					log.Fatalf("[ Redis][ AUTH] - %s\n", err.Error())
				}
			}
			if _, err = c.Do("SELECT", db); err != nil {
				c.Close()
				log.Fatalf("[ Redis][ SELECT] - redis(%s:%d) select db failed - %s",
					host, port, err.Error())
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

type IRedisStorage interface {
	// return a redis connections
	GetConn() redis.Conn
	// get keys start with prefix
	Keys(prefix string) ([]string, error)
}

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

func (r *redisStorage) checkOutputValueType(v interface{}) bool {
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

func (r *redisStorage) Source() interface{} {
	return r.pool
}

func (r *redisStorage) Driver() string {
	return DriveRedisStorage
}

func (r *redisStorage) decodeAssign(key string, dst interface{}) error {
	if r.checkOutputValueType(dst) {
		src, err := r.GetBytes(key)
		if err == nil {
			err = DecodeBytes(src, dst)
		}
		return err
	}
	return errors.New("dst must be struct")
}

func (r *redisStorage) Get(key string, dst interface{}) error {
	return r.decodeAssign(key, dst)
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
	conn := r.pool.Get()
	defer conn.Close()
	src, err := redis.String(conn.Do("GET", key))
	return src, err
}

func (r *redisStorage) GetBytes(key string) ([]byte, error) {
	conn := r.pool.Get()
	defer conn.Close()
	src, err := redis.Bytes(conn.Do("GET", key))
	return src, err
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
	return err == nil && i == 1
}

func (r *redisStorage) Delete(key string) {
	conn := r.pool.Get()
	defer conn.Close()
	conn.Do("DEL", key)
}

/*
   https://github.com/gomodule/redigo/issues/21
   It's common to store JSON in Redis. If you are only accessing
    the data from Go, then encoding/gob is another good option for storing nested data.
*/

func (r *redisStorage) hashSet(key string, v interface{}, seconds int) (err error) {
	if v == nil {
		return errors.New("nil value")
	}
	conn := r.pool.Get()
	defer conn.Close()
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Slice, reflect.Map, reflect.Struct:
		err = conn.Send("HMSET", redis.Args{key}.AddFlat(v)...)
		if seconds > 0 && err == nil {
			err = conn.Send("EXPIRE", key, seconds)
		}
		if err == nil {
			err = conn.Flush()
		}
		return err
	}
	err = conn.Send("SET", key, v)
	if seconds > 0 && err == nil {
		conn.Send("EXPIRE", key, seconds)
	}
	conn.Flush()
	return err
}

func (r *redisStorage) set(key string, v interface{}, seconds int64) error {
	conn := r.pool.Get()
	defer conn.Close()
	err := conn.Send("SET", key, v)
	if seconds > 0 && err == nil {
		conn.Send("EXPIRE", key, seconds)
	}
	conn.Flush()
	return err
}

func (r *redisStorage) binarySet(key string, v interface{}, seconds int64) error {
	byteData, err := EncodeBytes(v)
	if err == nil {
		return r.set(key, byteData, seconds)
	}
	return err
}

func (r *redisStorage) inputToBytes(v interface{}) bool {
	kind := reflect.TypeOf(v).Kind()
	return kind == reflect.Ptr || kind == reflect.Struct ||
		kind == reflect.Map || kind == reflect.Array
}

func (r *redisStorage) anySet(key string, v interface{}, seconds int64) error {
	if r.inputToBytes(v) {
		return r.binarySet(key, v, seconds)
	}
	return r.set(key, v, seconds)
}

func (r *redisStorage) Set(key string, v interface{}) error {
	return r.anySet(key, v, -1)
}

func (r *redisStorage) SetExpire(key string, v interface{}, seconds int64) error {
	if seconds > 0 {
		return r.anySet(key, v, seconds)
	}
	return nil
}

// return a redis connections
func (r *redisStorage) GetConn() redis.Conn {
	return r.pool.Get()
}

// get keys start with prefix
func (r *redisStorage) Keys(prefix string) ([]string, error) {
	conn := r.pool.Get()
	defer conn.Close()
	return redis.Strings(conn.Do("KEYS", prefix))
}

// delete keys contain prefix
func (r *redisStorage) DeleteWith(prefix string) (int, error) {
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

// Read and unmarshal from redis,if redis return err,
// marshal and write to redis
func (r *redisStorage) RWJson(key string, dst interface{},
	src func() interface{}, second int64) error {
	jsonBytes, err := r.GetBytes(key)
	if err == nil {
		err = json.Unmarshal(jsonBytes, &dst)
	}
	if err != nil {
		if src == nil {
			panic(errors.New("src is null pointer"))
		}
		dst = src()
		if dst != nil {
			jsonBytes, err = json.Marshal(dst)
			if err == nil {
				if second > 0 {
					r.SetExpire(key, jsonBytes, second)
				} else {
					r.Set(key, jsonBytes)
				}
			}
		}
	}
	return err
}
