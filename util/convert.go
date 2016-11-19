/**
 * Copyright 2015 @ at3.net.
 * name : convert.go
 * author : jarryliu
 * date : 2016-11-19 12:13
 * description :
 * history :
 */
package util

func I32Err(i int, err error) (int32, error) {
	if err != nil {
		return 0, err
	}
	return int32(i), err
}
