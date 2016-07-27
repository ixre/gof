/**
 * Copyright 2015 @ z3q.net.
 * name : storage.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */

package storage

// Storage
type Interface interface {
	// return storage driver
	Driver() interface{}

	// Return storage driver name
	DriverName() string
	// Check key is exists or not
	Exists(key string) (exists bool)
	// Get Value
	Get(key string, dst interface{}) error

	//Get raw value
	GetRaw(key string) (interface{}, error)

	// Set Value
	Set(key string, v interface{}) error

	GetBool(key string) (bool, error)

	GetInt(key string) (int, error)

	GetInt64(key string) (int64, error)

	GetString(key string) (string, error)

	GetFloat64(key string) (float64, error)

	// Delete Storage
	Del(key string)

	// Auto Delete Set
	SetExpire(key string, v interface{}, seconds int64) error
}
