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
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"math"
	mr "math/rand"
	"regexp"
	"time"
)

const letterStr = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var (
	revertRegexp = regexp.MustCompile("\\$\\{([^\\}]+)\\}")
)

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
func RandInt(n int) int {
	min := int(math.Pow10(n - 1))
	max := min*10 - 1
	mr.Seed(time.Now().Unix())
	v := mr.Intn(max)
	if v < min {
		return min + v
	}
	return v
}

//编码
func EncodingTransform(src []byte, enc string) ([]byte, error) {
	var ec encoding.Encoding
	switch enc {
	default:
		return src, nil
	case "GBK":
		ec = simplifiedchinese.GBK
	case "GB2312":
		ec = simplifiedchinese.HZGB2312
	case "BIG5":
		ec = traditionalchinese.Big5
	}
	dst := make([]byte, len(src)*2)
	n, _, err := ec.NewEncoder().Transform(dst, src, true)
	return dst[:n], err
}

// 还原模板的标签: ${...} -> {{...}}
func RevertTPVariable(str string) string {
	return revertRegexp.ReplaceAllString(str, "{{$1}}")
}
