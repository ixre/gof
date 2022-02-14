/**
 * Copyright 2015 @ 56x.net.
 * name : route.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package mvc

import (
	"github.com/ixre/gof/web"
	"log"
	"net/http"
	"strings"
)

var _ web.Route = new(Route)

type Route struct {
	_lazyRegister bool
	_routeMap     *web.RouteMap
	_ctlMap       map[string]ControllerGenerate
	_urlSuffix    string // 默认后缀
}

func NewRoute(source *web.RouteMap) *Route {
	if source == nil {
		source = new(web.RouteMap)
	}

	r := &Route{
		_ctlMap:    make(map[string]ControllerGenerate),
		_routeMap:  source,
		_urlSuffix: "",
	}
	return r
}

func (r *Route) SetSuffix(suffix string) {
	if suffix != "" {
		if suffix[0:1] != "." {
			suffix = "." + suffix
		}
		r._urlSuffix = suffix
	}
}

// Add 添加路由
func (r *Route) Add(urlPattern string, rf web.RequestHandler) {
	if urlPattern == "" || urlPattern == "*" {
		log.Fatalln("[ Panic] - Dangerous!The url parttern \"" + urlPattern +
			"\" will override default route," +
			"please instead of \"/\"")
	}
	r._routeMap.Add(urlPattern, rf)
}

// 处理请求
func (r *Route) Handle(ctx *web.Context) {
	if !r._lazyRegister {
		// 添加默认的路由
		r._routeMap.Add("*", r.handleAction)
		r._lazyRegister = true
	}
	r._routeMap.Handle(ctx)
}

// 延迟执行的操作，发生在请求完成后
func (r *Route) DeferFunc(f web.RequestHandler) {
	r._routeMap.DeferFunc(f)
}

func (r *Route) handleAction(ctx *web.Context) {

	path := ctx.Request.URL.Path
	var ctlName, action string
	ci := strings.Index(path[1:], "/")
	if ci == -1 {
		ctlName = path[1:]
	} else {
		ctlName = path[1 : ci+1]
	}

	action = GetAction(path, r._urlSuffix)
	if len(action) != 0 {
		if strings.ToLower(action) == ctlName {
			action = "Index"
		}

		if ctx.Request.Method == "POST" {
			action += "_post"
		}

		if r._ctlMap != nil {
			if v := r._ctlMap[ctlName]; v != nil {
				CustomHandle(v(), ctx, action)
				return
			}
		}
	}

	http.Error(ctx.Response, "404 page not found",
		http.StatusNotFound)
}

// 用普通方式注册控制器，由生成器为每个请求生成一个控制器实例。
func (r *Route) NormalRegister(name string, cg ControllerGenerate) {
	if r._ctlMap == nil {
		r._ctlMap = make(map[string]ControllerGenerate)
	}
	r._ctlMap[name] = cg
}

// 用单例方式注册控制器，所有的请求都共享一个控制器.
func (r *Route) Register(name string, c Controller) {
	var cg ControllerGenerate = func() Controller { return c }
	r.NormalRegister(name, cg)
}

// Get Controller
func (r *Route) GetController(name string) Controller {
	if r._ctlMap != nil {
		if v, ok := r._ctlMap[name]; ok {
			return v()
		}
	}
	return nil
}
