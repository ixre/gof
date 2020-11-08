package types

import "time"

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : time.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-07 20:42
 * description :
 * history :
 */

// parse chinese time text
func HanDateTime(t time.Time) string {
	return t.Format("2006年01月02日 15:04")
}
