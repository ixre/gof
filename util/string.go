/**
 * Copyright 2015 @ z3q.net.
 * name : string.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package util

import (
	"crypto/rand"
	mr "math/rand"
	"math"
)

const letterStr = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

//随机字符号串
func RandString(n int) string {
	lsLen := len(letterStr)
	var runes = make([]byte, n)
	rand.Read(runes)
	for i, b := range runes {
		runes[i] = letterStr[b%byte(lsLen)]
	}
	return string(runes)
}

// 随机整数
func RandInt(n int)int{
	min := int(math.Pow10(n-1))
	max := min * 10 -1
	v := mr.Intn(max)
	if v < min{
		return min + v
	}
	return v
}