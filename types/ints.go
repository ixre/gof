package types

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
