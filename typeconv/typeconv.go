package typeconv

import (
	"encoding/json"
	"fmt"
	"strconv"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : typeconv.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-21 20:19
 * description :
 * history :
 */

func parseInt(d interface{}, must bool) int {
	if d == nil && must {
		panic("parse nil to int fail")
	}
	switch d.(type) {
	case nil:
		return 0
	case string:
		i, err := strconv.Atoi(d.(string))
		if must && err != nil {
			panic("parse string to int fail:" + err.Error())
		}
		return i
	case float32:
		return int(d.(float32))
	case float64:
		return int(d.(float64))
	case int:
		return d.(int)
	case int8:
		return int(d.(int8))
	case int16:
		return int(d.(int16))
	case int32:
		return int(d.(int32))
	case int64:
		return int(d.(int64))
	}
	if must {
		panic("not support type:" + fmt.Sprintf("%+v", d))
	}
	return 0
}

// MustInt parse d to int
func Int(d interface{}) int {
	return parseInt(d, false)
}

// MustInt parse d to int
func MustInt(d interface{}) int {
	return parseInt(d, true)
}

// MustBool 将类型转换为bool
func MustBool(d interface{}) bool {
	switch d.(type) {
	case nil:
		return false
	case bool:
		return d.(bool)
	case int, int8, int16, int32, int64:
		return MustInt(d) == 1
	case string:
		b, _ := strconv.ParseBool(d.(string))
		return b
	}
	return false
}

// MustFloat parse d to float
func MustFloat(d interface{}) float64 {
	switch d.(type) {
	case nil:
		return 0
	case string:
		i, err := strconv.ParseFloat(d.(string), 64)
		if err != nil {
			panic("parse string to int fail:" + err.Error())
		}
		return i
	case float32:
		return float64(d.(float32))
	case float64:
		return d.(float64)
	case int:
		return float64(d.(int))
	case int8:
		return float64(d.(int8))
	case int16:
		return float64(d.(int16))
	case int32:
		return float64(d.(int32))
	case int64:
		return float64(d.(int64))
	}
	panic("not support type:" + fmt.Sprintf("%+v", d))
}

// Get string of interface, if can't converted,
// will return false
func String(d interface{}) (string, bool) {
	switch d.(type) {
	case string:
		return d.(string), true
	case []byte:
		return string(d.([]byte)), true
	case float32:
		return strconv.FormatFloat(float64(d.(float32)), 'g', 2, 32), true
	case float64:
		return strconv.FormatFloat(d.(float64), 'g', 2, 64), true
	case int:
		return strconv.FormatInt(int64(d.(int)), 10), true
	case int8:
		return strconv.FormatInt(int64(d.(int8)), 10), true
	case int16:
		return strconv.FormatInt(int64(d.(int16)), 10), true
	case int32:
		return strconv.FormatInt(int64(d.(int32)), 10), true
	case int64:
		return strconv.FormatInt(d.(int64), 10), true
	case uint:
		return strconv.FormatUint(uint64(d.(uint)), 10), true
	case uint8:
		return strconv.FormatUint(uint64(d.(uint8)), 10), true
	case uint16:
		return strconv.FormatUint(uint64(d.(uint16)), 10), true
	case uint32:
		return strconv.FormatUint(uint64(d.(uint32)), 10), true
	case uint64:
		return strconv.FormatUint(d.(uint64), 10), true
	case bool:
		return strconv.FormatBool(d.(bool)), true
	}
	return "", false
}

// get object string
func Stringify(d interface{}) string {
	if d == nil {
		return "null"
	}
	if s, b := String(d); b {
		return s
	}
	if d == nil {
		return "nil"
	}
	return fmt.Sprintf("%+v", d)
}

func MustJson(v interface{}) string {
	if v == nil {
		return "null"
	}
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func Int64Array(src []int64) []int {
	dst := make([]int, len(src))
	for i, v := range src {
		dst[i] = int(v)
	}
	return dst
}
