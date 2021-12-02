/**
 * Copyright 2015 @ 56x.net.
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
	"strings"
	"sync"
)

var _ Interface = new(hashStorage)

// 哈希表存储
type hashStorage struct {
	storage map[string]interface{}
	mux     *sync.RWMutex
}

func NewHashStorage() Interface {
	return &hashStorage{
		storage: make(map[string]interface{}),
		mux:     &sync.RWMutex{},
	}
}

// return storage driver
func (h *hashStorage) Source() interface{} {
	return h.storage
}

func (h *hashStorage) Driver() string {
	return DriveHashStorage
}

// check key is exists or not
func (h *hashStorage) Exists(key string) (exists bool) {
	_, b := h.test(key)
	return b
}

// test key,if key exists return value and true
func (h *hashStorage) test(key string) (interface{}, bool) {
	h.mux.RLock()
	v, ok := h.storage[key]
	h.mux.RUnlock()
	return v, ok
}

func (h *hashStorage) Get(key string, dst interface{}) error {
	k, ok := h.test(key)
	if ok {
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
	if v, err := h.GetRaw(key); err == nil && v != nil {
		if v2, ok := v.(bool); ok {
			return v2, nil
		}
	}
	return false, typeError
}

func (h *hashStorage) GetInt(key string) (int, error) {
	if v, err := h.GetRaw(key); err == nil && v != nil {
		if v2, ok := v.(int); ok {
			return v2, nil
		}
	}
	return 0, typeError
}

func (h *hashStorage) GetInt64(key string) (int64, error) {
	if v, err := h.GetRaw(key); err == nil && v != nil {
		if v2, ok := v.(int64); ok {
			return v2, nil
		}
	}
	return 0, typeError
}

func (h *hashStorage) GetString(key string) (string, error) {
	if v, err := h.GetRaw(key); err == nil && v != nil {
		if v2, ok := v.(string); ok {
			return v2, nil
		}
	}
	return "", typeError
}

func (h *hashStorage) GetBytes(key string) ([]byte, error) {
	if v, err := h.GetRaw(key); err == nil && v != nil {
		if v2, ok := v.([]byte); ok {
			return v2, nil
		}
	}
	return []byte(nil), typeError
}

func (h *hashStorage) GetFloat64(key string) (float64, error) {
	if v, err := h.GetRaw(key); err == nil && v != nil {
		if v2, ok := v.(float64); ok {
			return v2, nil
		}
	}
	return 0, typeError
}

func (h *hashStorage) Set(key string, v interface{}) error {
	h.mux.Lock()
	h.storage[key] = v
	h.mux.Unlock()
	return nil
}

//Get raw value
func (h *hashStorage) GetRaw(key string) (interface{}, error) {
	k, ok := h.test(key)
	if ok {
		return k, nil
	}
	return nil, errors.New("not such key")
}

func (h *hashStorage) Delete(key string) {
	h.mux.Lock()
	delete(h.storage, key)
	h.mux.Unlock()
}

func (h *hashStorage) DeleteWith(prefix string) (int, error) {
	h.mux.Lock()
	i := 0
	for k := range h.storage {
		if strings.HasPrefix(k, prefix) {
			delete(h.storage, k)
			i++
		}
	}
	h.mux.Unlock()
	return i, nil
}

// SetExpire equal of h.Set(key,value)
func (h *hashStorage) SetExpire(key string, v interface{}, seconds int64) error {
	return h.Set(key, v)
}

func (h *hashStorage) RWJson(key string, dst interface{}, src func() interface{}, second int64) error {
	err := h.Get(key, &dst)
	if err != nil {
		if src == nil {
			panic(errors.New("src is null pointer"))
		}
		dst = src()
		if dst != nil {
			h.SetExpire(key, dst, second)
		}
	}
	return err
}
