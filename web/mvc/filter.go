/**
 * Copyright 2015 @ z3q.net.
 * name : filter.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package mvc

import (
	"github.com/ixre/gof/web"
)

// controller filter
type Filter interface {
	//call it before execute your some business.
	Requesting(*web.Context) bool
	//call it after execute your some business.
	RequestEnd(*web.Context)
}
