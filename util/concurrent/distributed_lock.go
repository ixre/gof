package concurrent

import "github.com/ixre/gof/storage"

type DistributedLock struct {
	s storage.Interface
}

// NewDistributedLock 返回分步式锁
func NewDistributedLock(s storage.Interface) *DistributedLock {
	return &DistributedLock{
		s,
	}
}

// Lock 加锁, 返回true,表示加锁成功.否则已经加锁
func (d *DistributedLock) Lock(key string, expires int64) bool {
	k := d.getKey(key)
	if d.s.Exists(k) {
		return false
	}
	return d.s.SetExpire(k, "", expires) == nil
}

// Unlock 解锁
func (d *DistributedLock) Unlock(key string) {
	d.s.Delete(d.getKey(key))
}

func (d *DistributedLock) getKey(key string) string {
	return "_lock_" + key
}
