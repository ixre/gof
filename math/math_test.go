/**
 * Copyright 2015 @ at3.net.
 * name : math_test.go
 * author : jarryliu
 * date : 2016-10-01 12:25
 * description :
 * history :
 */
package math

import (
	"math"
	"strconv"
	"testing"
)

func TestNaN(t *testing.T) {
	f, err := strconv.ParseFloat("", 32)
	t.Log(f, err, f == math.NaN())
}
