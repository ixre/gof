/**
 * Copyright 2015 @ 56x.net.
 * name : storage.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */

package storage

import "errors"

const (
	DriveHashStorage  = "hash"
	DriveRedisStorage = "redis"
	DriveEtcdStorage  = "etcd"
)

var typeError = errors.New("type convert error")

// Storage
type Interface interface {
	// Return storage driver name
	Driver() string

	// return storage source
	Source() interface{}

	// check key is exists or not
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
	Delete(key string)
	// delete by key prefix
	DeleteWith(prefix string) (int, error)
	// Read and unmarshal from redis,if redis return err,
	// marshal and write to redis
	RWJson(key string, dst interface{}, src func() interface{}, second int64) error
}
