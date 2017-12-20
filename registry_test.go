package gof

import (
	"testing"
)

func TestRegistry_Set(t *testing.T) {
	r, _ := NewRegistry("./tmp/conf/", ".")
	rt := r.Use("core")
	key := "config.user_name"
	val := rt.Get(key)
	if val == nil {
		err := rt.Set(key, "jarrysix")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		rt.Flush()
		val = rt.Get(key)
	}
	t.Log("result :", val.(string))
}

func TestRegistryOld(t *testing.T) {
	r, _ := NewRegistry("./tmp/conf/", ".")
	key := "core.config.user_name"
	val := r.Get(key)
	t.Log("result :", val.(string))
}
