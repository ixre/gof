package types

import (
	"fmt"
	"strconv"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : strings.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-09-13 10:28
 * description :
 * history :
 */

// Get string of interface, if can't converted,
// will return false
func ToString(d interface{}) (string, bool) {
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

func String(d interface{}) string {
	s, b := ToString(d)
	if !b {
		s = fmt.Sprintf("%+v", d)
	}
	return s
}
