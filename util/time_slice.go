package util

import (
	"strconv"
	"strings"
	"time"
)

/**
 * Copyright 2009-2019 @ 56x.net
 * name : time_ticker
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-10-03 16:18
 * description :
 * history :
 */

// 获取刻度字符串
func getTickString(t time.Time, unitD, unitH, unitM, unitS int) string {
	var d, h, m, s int
	if unitD > 0 {
		d = t.Day() / unitD
	} else if unitH > 0 {
		h = t.Hour() / unitH
	} else if unitM > 0 {
		m = t.Minute() / unitM
	} else if unitS > 0 {
		s = t.Second() / unitS
	}
	return strings.Join([]string{
		strconv.Itoa(d),
		strconv.Itoa(h),
		strconv.Itoa(m),
		strconv.Itoa(s),
	}, "-")
}

// 获取时间(小时)切片字符串, 如: unit = 15, t = 12:00:14 计算结果可能为: 01-00-00
func GetHourSlice(t time.Time, unit int) string {
	return getTickString(t, 0, unit, 0, 0)
}

// 获取时间(分钟)切片字符串, 如: unit = 15, t = 12:00:14 计算结果可能为: 00-01-00
func GetMinuteSlice(t time.Time, unit int) string {
	return getTickString(t, 0, 0, unit, 0)
}

// 获取时间(秒)切片字符串, 如: unit = 15, t = 12:00:14 计算结果可能为: 00-00-01
func GetSecondSlice(t time.Time, unit int) string {
	return getTickString(t, 0, 0, 0, unit)
}
