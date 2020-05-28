package util

import (
	"fmt"
	"strconv"
)

/**
 * Copyright 2009-2019 @ to2.net
 * name : types
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-11-18 10:17
 * description :
 * history :
 */

func Str(d interface{}) string {
	return TypeToString(d)
}

// 将类型转为string
func TypeToString(d interface{}) string {
	switch d.(type) {
	case string:
		return d.(string)
	case []byte:
		return string(d.([]byte))
	case float32:
		return strconv.FormatFloat(float64(d.(float32)), 'g', 2, 32)
	case float64:
		return strconv.FormatFloat(d.(float64), 'g', 2, 64)
	case int:
		return strconv.FormatInt(int64(d.(int)), 10)
	case int8:
		return strconv.FormatInt(int64(d.(int8)), 10)
	case int16:
		return strconv.FormatInt(int64(d.(int16)), 10)
	case int32:
		return strconv.FormatInt(int64(d.(int32)), 10)
	case int64:
		return strconv.FormatInt(d.(int64), 10)
	case uint:
		return strconv.FormatUint(uint64(d.(uint)), 10)
	case uint8:
		return strconv.FormatUint(uint64(d.(uint8)), 10)
	case uint16:
		return strconv.FormatUint(uint64(d.(uint16)), 10)
	case uint32:
		return strconv.FormatUint(uint64(d.(uint32)), 10)
	case uint64:
		return strconv.FormatUint(d.(uint64), 10)
	case bool:
		return strconv.FormatBool(d.(bool))
	}
	return fmt.Sprintf("%+v", d)
}

// 将类型转为string
func TypeToInt(d interface{}) int {
	switch d.(type) {
	case string:
		i, err := strconv.Atoi(d.(string))
		if err != nil {
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
		return d.(int)
	case int16:
		return d.(int)
	case int32:
		return d.(int)
	case int64:
		return d.(int)
	}
	panic("not support type:" + fmt.Sprintf("%+v", d))
}
