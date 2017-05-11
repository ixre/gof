/**
 * Copyright 2015 @ at3.net.
 * name : storage_test.go
 * author : jarryliu
 * date : 2016-11-24 21:23
 * description :
 * history :
 */
package storage

import "testing"

type testClass struct {
	Age  int
	Name string
}

func getRdsStorage() Interface {
	pool := NewRedisPool("dbs.ts.com", 6379, 1,
		"123456", 0, 0)
	return NewRedisStorage(pool)
}

func TestNewRedisPool(t *testing.T) {
	k1 := "key:test:hash"
	rds := getRdsStorage()
	c := &testClass{
		Age:  10,
		Name: "100",
	}
	rds.SetExpire(k1, c, 3600)
	c.Age = 20
	c.Name = "200"
	err := rds.Get(k1, c)
	t.Log(c.Age, c.Name, err)

	i := 1
	k2 := "key:test:int2"
	rds.Set(k2, i)
	i = 2
	i, _ = rds.GetInt(k2)
	t.Log(i)

}
