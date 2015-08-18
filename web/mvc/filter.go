/**
 * Copyright 2015 @ S1N1 Team.
 * name : filter.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package mvc

import (
	"github.com/jrsix/gof/web"
)

// controller filter
type Filter interface {
	//call it before execute your some business.
	Requesting(*web.Context) bool
	//call it after execute your some business.
	RequestEnd(*web.Context)
}
