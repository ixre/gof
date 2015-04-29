/**
 * Copyright 2015 @ S1N1 Team.
 * name : storage.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */

package gof

type Storage interface {
	// Get Value
	Get(key string, dst interface{}) error
	// Set Value
	Set(key string, v interface{})error
	// Auto Delete Set
	DSet(key string, v interface{}, seconds int32)
}
