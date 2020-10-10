package api

import (
	"testing"
)

func TestVersionCompare(t *testing.T) {
	v1, v2 := "1.3.1.100", "1.3.1.100"
	b := CompareVersion(v1, v2)
	b2 := IntVersion(v2) > IntVersion(v1)
	t.Log(v1, ">", v2, " | ", b, b2)
}

type MyStruct struct {
}

func TestLang(t *testing.T) {
	arr := []*MyStruct{{}}
	do(arr)
	t.Log(arr[0] == nil)
}
func do(arr []*MyStruct) {
	arr[0] = nil
}
