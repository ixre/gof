/**
 * Copyright 2015 @ 56x.net.
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

// Filter controller filter
type Filter interface {
	// Requesting call it before execute your some business.
	Requesting(*web.Context) bool
	// RequestEnd call it after execute your some business.
	RequestEnd(*web.Context)
}
