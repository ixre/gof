/**
 * Copyright 2015 @ z3q.net.
 * name : storage.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */

package storage

import "errors"

const (
	DriveHashStorage  = "hash-storage"
	DriveRedisStorage = "redis-storage"
)

var typeError = errors.New("type convert error")

// Storage
type Interface interface {
	// Return storage driver name
	Driver() string

	// return storage source
	Source() interface{}

	// Check key is exists or not
	Exists(key string) (exists bool)

	// Set Value
	Set(key string, v interface{}) error

	// Auto Delete Set
	SetExpire(key string, v interface{}, seconds int64) error

	// Get Value
	Get(key string, dst interface{}) error

	//Get raw value
	GetRaw(key string) (interface{}, error)

	GetBool(key string) (bool, error)

	GetInt(key string) (int, error)

	GetInt64(key string) (int64, error)

	GetString(key string) (string, error)

	GetFloat64(key string) (float64, error)

	GetBytes(key string) ([]byte, error)

	// Delete Storage
	Del(key string)

	// Read and unmarshal from redis,if redis return err,
	// marshal and write to redis
	RWJson(key string, dst interface{}, src func() interface{}, second int64) error
}
