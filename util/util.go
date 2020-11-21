package util

import (
	"github.com/ixre/gof/types"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var LetterBytes = []byte("01234ABCDEFGHIJK56789abcdefghijklmLMNOPQRSTUVWXYZnopqrstuvwxyz")

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

func (e *stringExtend) String(v interface{}) string {
	return types.Stringify(v)
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
