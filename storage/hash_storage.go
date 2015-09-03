/**
 * Copyright 2015 @ z3q.net.
 * name : map_storage
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package storage

import (
	"errors"
	"github.com/jsix/gof"
	"sync"
)

var DriveHashStorage string = "hash-storage"
var typeError error = errors.New("type convert error!")

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

func (this *hashStorage) GetBool(key string) (bool, error) {
	if v, _ := this.GetRaw(key); v != nil {
		if v2, ok := v.(bool); ok {
			return v2, nil
		}
	}
	return false, typeError
}

func (this *hashStorage) GetInt(key string) (int, error) {
	if v, _ := this.GetRaw(key); v != nil {
		if v2, ok := v.(int); ok {
			return v2, nil
		}
	}
	return 0, typeError
}

func (this *hashStorage) GetInt64(key string) (int64, error) {
	if v, _ := this.GetRaw(key); v != nil {
		if v2, ok := v.(int64); ok {
			return v2, nil
		}
	}
	return 0, typeError
}

func (this *hashStorage) GetString(key string) (string, error) {
	if v, _ := this.GetRaw(key); v != nil {
		if v2, ok := v.(string); ok {
			return v2, nil
		}
	}
	return "", typeError
}

func (this *hashStorage) GetFloat64(key string) (float64, error) {
	if v, _ := this.GetRaw(key); v != nil {
		if v2, ok := v.(float64); ok {
			return v2, nil
		}
	}
	return 0, typeError
}

func (this *hashStorage) Set(key string, v interface{}) error {
	this._map[key] = v
	return nil
}

//Get raw value
func (this *hashStorage) GetRaw(key string) (interface{}, error) {
	if k, ok := this._map[key]; ok {
		return k, nil
	}
	return nil, errors.New("not such key")
}

func (this *hashStorage) Del(key string) {
	delete(this._map, key)
}

func (this *hashStorage) SetExpire(key string, v interface{}, seconds int64) error {
	panic(errors.New("HashStorage not support method \"SetExpire\"!"))
}
