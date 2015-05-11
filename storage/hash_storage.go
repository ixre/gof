/**
 * Copyright 2015 @ S1N1 Team.
 * name : map_storage
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package storage

import (
	"errors"
	"github.com/atnet/gof"
	"sync"
)

var DriveHashStorage string = "hash-storage"

// 哈希表存储
type hashStorage struct {
	_map map[string]interface{}
	sync.Mutex
}

func NewHashStorage() gof.Storage {
	return &hashStorage{
		_map: make(map[string]interface{}),
	}
}

func (this *hashStorage) Driver() string {
	return DriveHashStorage
}

func (this *hashStorage) Get(key string, dst interface{}) error {
	panic(errors.New("HashStorage not support method \"Get\"!"))
}

func (this *hashStorage) Set(key string, v interface{}) error {
	this._map[key] = v
	return nil
}

//Get raw value
func (this *hashStorage) GetRaw(key string) interface{} {
	if k, ok := this._map[key]; ok {
		return k
	}
	return nil
}

func (this *hashStorage) Del(key string) {
	delete(this._map, key)
}

func (this *hashStorage) SetExpire(key string, v interface{}, seconds int64) error {
	panic(errors.New("HashStorage not support method \"SetExpire\"!"))
}
