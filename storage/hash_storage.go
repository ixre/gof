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
	"reflect"
	"sync"
)

var _ Interface = new(hashStorage)

// 哈希表存储
type hashStorage struct {
	storage map[string]interface{}
	sync.Mutex
}

func NewHashStorage() Interface {
	return &hashStorage{
		storage: make(map[string]interface{}),
	}
}

// return storage driver
func (h *hashStorage) Source() interface{} {
	return h.storage
}

func (h *hashStorage) Driver() string {
	return DriveHashStorage
}

// Check key is exists or not
func (h *hashStorage) Exists(key string) (exists bool) {
	_, b := h.storage[key]
	return b
}

func (h *hashStorage) Get(key string, dst interface{}) error {
	if k, ok := h.storage[key]; ok {
		if reflect.TypeOf(k).Kind() == reflect.Ptr {
			dst = k
		} else {
			dst = &k
		}
		return nil
	}
	return errors.New("not such key")
}

func (h *hashStorage) GetBool(key string) (bool, error) {
	if v, _ := h.GetRaw(key); v != nil {
		if v2, ok := v.(bool); ok {
			return v2, nil
		}
	}
	return false, typeError
}

func (h *hashStorage) GetInt(key string) (int, error) {
	if v, _ := h.GetRaw(key); v != nil {
		if v2, ok := v.(int); ok {
			return v2, nil
		}
	}
	return 0, typeError
}

func (h *hashStorage) GetInt64(key string) (int64, error) {
	if v, _ := h.GetRaw(key); v != nil {
		if v2, ok := v.(int64); ok {
			return v2, nil
		}
	}
	return 0, typeError
}

func (h *hashStorage) GetString(key string) (string, error) {
	if v, _ := h.GetRaw(key); v != nil {
		if v2, ok := v.(string); ok {
			return v2, nil
		}
	}
	return "", typeError
}

func (h *hashStorage) GetFloat64(key string) (float64, error) {
	if v, _ := h.GetRaw(key); v != nil {
		if v2, ok := v.(float64); ok {
			return v2, nil
		}
	}
	return 0, typeError
}

func (h *hashStorage) Set(key string, v interface{}) error {
	h.storage[key] = v
	return nil
}

//Get raw value
func (h *hashStorage) GetRaw(key string) (interface{}, error) {
	if k, ok := h.storage[key]; ok {
		return k, nil
	}
	return nil, errors.New("not such key")
}

func (h *hashStorage) Del(key string) {
	delete(h.storage, key)
}

func (h *hashStorage) SetExpire(key string, v interface{}, seconds int64) error {
	panic(errors.New("HashStorage not support method \"SetExpire\"!"))
}
