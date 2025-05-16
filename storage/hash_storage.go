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
	"time"
)

var _ Interface = new(hashStorage)

// 存储项，包含值和过期时间
type storageItem struct {
	value     interface{}
	expiresAt int64
}

// 哈希表存储
type hashStorage struct {
	storage map[string]storageItem
	mux     *sync.RWMutex
}

func NewHashStorage() Interface {
	hs := &hashStorage{
		storage: make(map[string]storageItem),
		mux:     &sync.RWMutex{},
	}
	// 启动清理协程
	go hs.startCleaner()
	return hs
}

// return storage driver
func (h *hashStorage) Source() interface{} {
	return h.storage
}

func (h *hashStorage) Driver() string {
	return DriveHashStorage
}

// 检查键是否存在且未过期
func (h *hashStorage) Exists(key string) (exists bool) {
	_, b := h.test(key)
	return b
}

// 测试键，如果键存在且未过期返回值和 true
func (h *hashStorage) test(key string) (interface{}, bool) {
	h.mux.RLock()
	defer h.mux.RUnlock()
	item, ok := h.storage[key]
	if !ok {
		return nil, false
	}
	if item.expiresAt != 0 &&
		time.Now().Unix() > item.expiresAt {
		return nil, false
	}
	return item.value, true
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
	return h.SetExpire(key, v, 0)
}

// Set 设置值并可指定过期时间，seconds 为 0 表示永不过期
func (h *hashStorage) SetExpire(key string, v interface{}, seconds int64) error {
	h.mux.Lock()
	defer h.mux.Unlock()
	var expiresAt int64
	if seconds > 0 {
		expiresAt = time.Now().Unix() + seconds
	}
	h.storage[key] = storageItem{
		value:     v,
		expiresAt: expiresAt,
	}
	return nil
}

// Get raw value
func (h *hashStorage) GetRaw(key string) (interface{}, error) {
	k, ok := h.test(key)
	if ok {
		return k, nil
	}
	return nil, errors.New("not such key")
}

func (h *hashStorage) Delete(key string) {
	h.mux.Lock()
	defer h.mux.Unlock()
	delete(h.storage, key)
}

func (h *hashStorage) DeleteWith(prefix string) (int, error) {
	h.mux.Lock()
	defer h.mux.Unlock()
	i := 0
	for k := range h.storage {
		if strings.HasPrefix(k, prefix) {
			delete(h.storage, k)
			i++
		}
	}
	return i, nil
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

// 添加定期清理方法
func (h *hashStorage) startCleaner() {
	ticker := time.NewTicker(time.Second * 10) // 每10秒清理一次
	defer ticker.Stop()

	for range ticker.C {
		h.cleanExpired()
	}
}

// 添加清理过期项的方法
func (h *hashStorage) cleanExpired() {
	h.mux.Lock()
	defer h.mux.Unlock()

	now := time.Now().Unix()
	for k, v := range h.storage {
		if v.expiresAt > 0 && v.expiresAt < now {
			delete(h.storage, k)
		}
	}
}
