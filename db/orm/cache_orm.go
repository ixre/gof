/**
 * Copyright 2015 @ z3q.net.
 * name : cache_orm.go
 * author : jarryliu
 * date : 2016-07-27 11:48
 * description :
 * history :
 */
package orm

import (
	"github.com/jsix/gof/storage"
)

var _ Orm = new(cacheProxy)

type cacheProxy struct {
	Orm
	Storage storage.Interface
}

func CacheProxy(o Orm, s storage.Interface) {
	o = NewCacheProxyOrm(o, s)
}

func NewCacheProxyOrm(o Orm, s storage.Interface) Orm {
	return &cacheProxy{
		Orm:     o,
		Storage: s,
	}
}
