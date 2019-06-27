package types

/**
 * Copyright 2009-2019 @ to2.net
 * name : types_.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-06-26 22:15
 * description :
 * history :
 */


func StringDefault(s,d string)string{
	if len(s)==0 {
		return d
	}
	return s
}