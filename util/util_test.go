package util

import (
	"testing"
)

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

func TestMergeIntArray(t *testing.T) {
	old := []int{1,2,3,4,5,6}
	n := []int{3,2,8,9}
	final,del := IntArrayDiff(old,n,nil)
	t.Log(final)
	t.Log(del)
}
