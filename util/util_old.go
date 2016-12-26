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
	"encoding/json"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"html/template"
	"log"
	"math"
	mr "math/rand"
	"time"
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

// 强制序列化为可用于HTML的JSON
func MustHtmlJson(v interface{}) template.JS {
	d, err := json.Marshal(v)
	if err != nil {
		log.Println("[ Json][ Mashal]: ", err.Error())
	}
	return template.JS(d)
}
