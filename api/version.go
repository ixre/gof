package api

import (
	"strconv"
	"strings"
)

func CompareVersion(v, v1 string) int {
	return IntVersion(v) - IntVersion(v1)
}
func IntVersion(s string) int {
	arr := strings.Split(s, ".")
	for i, v := range arr {
		if l := len(v); l < 3 {
			arr[i] = strings.Repeat("0", 3-l) + v
		}
	}
	intVer, err := strconv.Atoi(strings.Join(arr, ""))
	if err != nil {
		panic(err)
	}
	return intVer
}

