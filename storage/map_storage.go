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
    "github.com/atnet/gof"
    "sync"
    "errors"
)

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

func (this *hashStorage) Get(key string, dst interface{}) error {
    if k, ok := this._map[key]; ok {
        dst = k
    }
    return nil
}

func (this *hashStorage) Set(key string, v interface{}) error {
    this._map[key] = v
    return nil
}

func (this *hashStorage) Del(key string) {
    delete(this._map,key)
}

func (this *hashStorage) SetExpire(key string, v interface{}, seconds int64) error {
    return errors.New("HashStorage not support method \"SetExpire\"!")
}
