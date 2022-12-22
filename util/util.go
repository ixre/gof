package util

import (
	"encoding/json"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/ixre/gof/types/typeconv"
)

var LetterBytes = []byte("01234ABCDEFGHIJK56789abcdefghijklmLMNOPQRSTUVWXYZnopqrstuvwxyz")

// RandString 随机字符号串
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

// RandInt 随机整数
func RandInt(n int) int {
	min := int(math.Pow10(n - 1))
	max := min*10 - 1
	rand.Seed(time.Now().UnixNano())
	v := rand.Intn(max)
	if v < min {
		return min + v
	}
	return v
}



// Stringify 返回对象的JSON格式数据
func Stringify(o interface{}) string {
	bytes, _ := json.Marshal(o)
	return string(bytes)
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

func (e *stringExtend) String(v interface{}) string {
	return typeconv.Stringify(v)
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

// 数组拼接字符串
func JoinIntArray(arr []int, sep string) string {
	var strIds = make([]string, len(arr))
	for i, v := range arr {
		strIds[i] = strconv.Itoa(v)
	}
	return strings.Join(strIds, sep)
}

// 比较数组差异
func IntArrayDiff(o []int, n []int, fn func(v int, add bool)) (created []int, deleted []int) {
	isExists := func(arr []int, v int) bool {
		for _, x := range arr {
			if x == v {
				return true
			}
		}
		return false
	}
	exists := []int{}
	// 旧数组查找已删除
	for _, v := range o {
		if !isExists(n, v) {
			deleted = append(deleted, v)
		} else {
			exists = append(exists, v)
		}
	}
	// 查找新增
	for _, v := range n {
		if !isExists(o, v) {
			created = append(created, v)
		}
	}
	if fn != nil {
		for _, v := range exists {
			fn(v, false)
		}
		for _, v := range created {
			fn(v, true)
		}
	}
	return created, deleted
}
