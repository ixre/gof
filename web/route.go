/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : newmin
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package web

import (
	"net/http"
	"regexp"
)

//路由映射
type RouteMap struct {
	deferFunc HttpContextFunc
	//地址模式
	UrlPatterns []string
	//路由集合
	RouteCollection map[string]HttpContextFunc
}

// HTTP处理词典
type httpFuncMap map[string]HttpContextFunc

//添加路由
func (this *RouteMap) Add(urlPattern string,rf HttpContextFunc) {
	if this.RouteCollection == nil {
		this.RouteCollection = make(map[string]HttpContextFunc)
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
	r,w := ctx.Request,ctx.ResponseWriter
	routes := this.RouteCollection
	path := r.URL.Path
	var isHandled bool = false

	//range 顺序是随机的，参见：http://golanghome.com/post/155
	for _, k := range this.UrlPatterns {
		v, exist := routes[k]
		if exist {
			matched, err := regexp.Match(k, []byte(path))
			//fmt.Println("\n",k,path)
			if err != nil {
				http.Error(w,err.Error(),http.StatusInternalServerError)
				return
			}
			if matched && v != nil {
				isHandled = true
				v(ctx)
				break
			}
		}
	}

	if !isHandled {
		http.Error(w,"404 Not found!",http.StatusNotFound)
	}
}

// 延迟执行的操作，发生在请求完成后
func (this *RouteMap) Defer(f HttpContextFunc) {
	this.deferFunc = f
}