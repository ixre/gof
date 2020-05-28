package util

import "testing"

func TestRandomString(t *testing.T) {
	mp := make(map[string]int)
	for i := 0; i < 100000; i++ {
		str := RandString(6)
		t.Logf("%s\n", str)
		if _, ok := mp[str]; ok {
			t.Log("重复")
			t.FailNow()
		} else {
			mp[str] = 1
		}
	}
}
