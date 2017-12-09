package gof

import "testing"

func TestRegistry_Set(t *testing.T) {
	r,err := NewRegistry("./tmp/conf/", ":")
	key := "core:config:user_name"
	err = r.Set(key, "jarrysix")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	val := r.Get(key)
	t.Log("result :", val.(string))
}
