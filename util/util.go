package util

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var LetterBytes = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

//随机字符号串
func RandString(n int) string {
	return string(RandBytes(n, LetterBytes))
}

func RandBytes(n int, letters []byte) []byte {
	l := len(letters)
	var arr = make([]byte, n)
	rand.Read(arr)
	for i, b := range arr {
		arr[i] = letters[b%byte(l)]
	}
	return arr
}

// 随机整数
func RandInt(n int) int {
	min := int(math.Pow10(n - 1))
	max := min*10 - 1
	rand.Seed(time.Now().Unix())
	v := rand.Intn(max)
	if v < min {
		return min + v
	}
	return v
}

var (
	BoolExt = &boolExtend{}
	StrExt  = &stringExtend{}
)

type (
	boolExtend struct {
	}
	stringExtend struct {
	}
)

func threeCondition(b bool, i1, i2 interface{}) interface{} {
	if b {
		return i1
	}
	return i2
}

func (e *boolExtend) TInt(b bool, v1, v2 int) int {
	return threeCondition(b, v1, v2).(int)
}
func (e *boolExtend) TInt32(b bool, v1, v2 int32) int32 {
	return threeCondition(b, v1, v2).(int32)
}
func (e *boolExtend) TInt64(b bool, v1, v2 int64) int64 {
	return threeCondition(b, v1, v2).(int64)
}
func (e *boolExtend) TString(b bool, v1, v2 string) string {
	return threeCondition(b, v1, v2).(string)
}

func (e *stringExtend) String(v interface{}) string {
	return Str(v)
}

// 字符串转为int32切片
func (e *stringExtend) I32Slice(s string, delimer string) []int32 {
	var arr []int32
	sArr := strings.Split(s, delimer)
	for _, v := range sArr {
		i, err := strconv.ParseInt(v, 10, 32)
		if err == nil {
			arr = append(arr, int32(i))
		}
	}
	return arr
}

// 字符串转为int切片
func (e *stringExtend) IntSlice(s string, delimer string) []int {
	var arr []int
	sArr := strings.Split(s, delimer)
	for _, v := range sArr {
		i, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			arr = append(arr, int(i))
		}
	}
	return arr
}
