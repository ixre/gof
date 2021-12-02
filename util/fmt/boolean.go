/**
 * Copyright 2015 @ 56x.net.
 * name : fmt
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */

package fmt

func BoolString(b bool, trueVal, falseVal string) string {
	if b {
		return trueVal
	}
	return falseVal
}

func BoolInt(b bool, v, v1 int) int {
	if b {
		return v
	}
	return v1
}
