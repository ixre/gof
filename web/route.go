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
	"errors"
	"net/http"
	"regexp"
	"strings"
)

// Url Route
type Route interface {
	Add(urlPattern string, rf RequestHandler)
	Handle(ctx *Context)
}

var _ Route = new(RouteMap)

//路由映射
type RouteMap struct {
	deferFunc RequestHandler
	//地址模式
	UrlPatterns []string
	//路由集合
	RouteCollection map[string]RequestHandler
}

// HTTP处理词典
type httpFuncMap map[string]RequestHandler

//添加路由
func (this *RouteMap) Add(urlPattern string, rf RequestHandler) {
	if this.RouteCollection == nil {
		this.RouteCollection = make(map[string]RequestHandler)
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
	r, w := ctx.Request, ctx.Response
	routes := this.RouteCollection
	path := r.URL.Path
	var isHandled bool = false

	//range 顺序是随机的，参见：http://golanghome.com/post/155
	for _, routeKey := range this.UrlPatterns {
		if routeHandler, exist := routes[routeKey]; exist {
			if isHandled, err = this.chkInvoke(path, routeKey, routeHandler, ctx); isHandled {
				break
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				break
			}
		}
	}

	if !isHandled {
		http.Error(w, "404 Not found!", http.StatusNotFound)
	}
}

func (this *RouteMap) chkInvoke(requestPath, routeKey string, routeHandler RequestHandler, ctx *Context) (bool, error) {
	if routeHandler == nil {
		panic(errors.New("handler can't nil!"))
	}

	// 是否为通配的路由(*)，或与路由规则一致
	if match := routeKey == "*" || routeKey == requestPath; match {
		return true, this.callHandler(routeHandler, ctx)
	}

	// 如果路由不符，且规则为正则，前尝试匹配
	if routeKey[0:1] == "^" {
		if match, err := regexp.MatchString(routeKey, requestPath); match {
			return true, this.callHandler(routeHandler, ctx)
		} else {
			return false, err
		}
	}

	// 如果结尾为“*”，标题匹配以前的URL
	var j int = len(routeKey) - 1
	if routeKey[j:] == "*" {
		if strings.HasPrefix(requestPath, routeKey[:j]) {
			return true, this.callHandler(routeHandler, ctx)
		}
	}

	return false, nil
}

func (this *RouteMap) callHandler(handler RequestHandler, ctx *Context) error {
	if handler == nil {
		return errors.New("No handler process your request!")
	}
	handler(ctx)
	return nil
}

// 延迟执行的操作，发生在请求完成后
func (this *RouteMap) DeferFunc(f RequestHandler) {
	this.deferFunc = f
}
