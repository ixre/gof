package types

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : types.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-09-19 14:01
 * description :
 * history :
 */

func threeCondition(b bool, i1, i2 interface{}) interface{} {
	if b {
		return i1
	}
	return i2
}

// Int条件判断
func IntCond(b bool, v1, v2 int) int {
	return threeCondition(b, v1, v2).(int)
}
func StringCond(b bool, v1, v2 string) string {
	return threeCondition(b, v1, v2).(string)
}
func FloatCond(b bool, v1, v2 float64) float64 {
	return threeCondition(b, v1, v2).(float64)
}

func TInt32(b bool, v1, v2 int32) int32 {
	return threeCondition(b, v1, v2).(int32)
}
func TInt64(b bool, v1, v2 int64) int64 {
	return threeCondition(b, v1, v2).(int64)
}
