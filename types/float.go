package types

import (
	"fmt"
	"strings"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : float.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-07 20:43
 * description :
 * history :
 */

// parse to money string like '0.60'
func FixedMoney(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

// parse to money string, but defence by FixedMoney, it's like '0.6'
func Money(f float64) string {
	s := FixedMoney(f)
	if strings.HasSuffix(s, ".00") {
		return s[:len(s)-3]
	} else if strings.HasSuffix(s, "0") {
		return s[:len(s)-1]
	}
	return s
}
