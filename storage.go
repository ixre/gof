/**
 * Copyright 2015 @ S1N1 Team.
 * name : storage.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */

package gof

// Storage
type Storage interface {
	// Return storage drive name
	Driver() string

	// Get Value
	Get(key string, dst interface{}) error

	//Get raw value
	GetRaw(key string) interface{}

	// Set Value
	Set(key string, v interface{}) error

	// Delete Storage
	Del(key string)

	// Auto Delete Set
	SetExpire(key string, v interface{}, seconds int64) error
}
