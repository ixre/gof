/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package web

import (
	"net/http"
	"regexp"
)

// Url Route
type Route interface {
	Add(urlPattern string, rf ContextFunc)
	Handle(ctx *Context)
}

var _ Route = new(RouteMap)

//路由映射
type RouteMap struct {
	deferFunc ContextFunc
	//地址模式
	UrlPatterns []string
	//路由集合
	RouteCollection map[string]ContextFunc
}

// HTTP处理词典
type httpFuncMap map[string]ContextFunc

//添加路由
func (this *RouteMap) Add(urlPattern string, rf ContextFunc) {
	if this.RouteCollection == nil {
		this.RouteCollection = make(map[string]ContextFunc)
	}
	_, exists := this.RouteCollection[urlPattern]
	if !exists {
		this.RouteCollection[urlPattern] = rf
		this.UrlPatterns = append(this.UrlPatterns, urlPattern)
	}
}

// 处理请求
func (this *RouteMap) Handle(ctx *Context) {
	// 执行某些操作，如捕获异常等
	if this.deferFunc != nil {
		defer this.deferFunc(ctx)
	}
	var err error
	r, w := ctx.Request, ctx.ResponseWriter
	routes := this.RouteCollection
	path := r.URL.Path
	var isHandled bool = false

	//range 顺序是随机的，参见：http://golanghome.com/post/155
	for _, k := range this.UrlPatterns {
		v, exist := routes[k]
		if exist {
			var isMatch bool

			// 是否为通配的路由(*)，或与路由规则一致
			isMatch = k == "*" || k == path

			// 如果路由不符，且规则为正则，前尝试匹配
			if !isMatch && k[0:1] == "^" {
				isMatch, err = regexp.MatchString(k, path)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			//fmt.Println("Verify:", k, path)
			if isMatch && v != nil {
				//fmt.Println("Matched:", k, path)
				isHandled = true
				v(ctx)
				break
			}
		}
	}

	if !isHandled {
		http.Error(w, "404 Not found!", http.StatusNotFound)
	}
}

// 延迟执行的操作，发生在请求完成后
func (this *RouteMap) DeferFunc(f ContextFunc) {
	this.deferFunc = f
}
