/**
 * Copyright 2015 @ at3.net.
 * name : convert.go
 * author : jarryliu
 * date : 2016-11-19 12:13
 * description :
 * history :
 */
package util

import (
    "strconv"
    "fmt"
)

func I32Err(i int, err error) (int32, error) {
    if err != nil {
        return 0, err
    }
    return int32(i), err
}

// 将类型转为string
func Str(d interface{}) string {
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