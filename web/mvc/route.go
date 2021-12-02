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

func (this *Route) SetSuffix(suffix string) {
	if suffix != "" {
		if suffix[0:1] != "." {
			suffix = "." + suffix
		}
		this._urlSuffix = suffix
	}
}

//添加路由
func (this *Route) Add(urlPattern string, rf web.RequestHandler) {
	if urlPattern == "" || urlPattern == "*" {
		log.Fatalln("[ Panic] - Dangerous!The url parttern \"" + urlPattern +
			"\" will override default route," +
			"please instead of \"/\"")
	}
	this._routeMap.Add(urlPattern, rf)
}

// 处理请求
func (this *Route) Handle(ctx *web.Context) {
	if !this._lazyRegister {
		// 添加默认的路由
		this._routeMap.Add("*", this.handleAction)
		this._lazyRegister = true
	}
	this._routeMap.Handle(ctx)
}

// 延迟执行的操作，发生在请求完成后
func (this *Route) DeferFunc(f web.RequestHandler) {
	this._routeMap.DeferFunc(f)
}

func (this *Route) handleAction(ctx *web.Context) {

	path := ctx.Request.URL.Path
	var ctlName, action string
	ci := strings.Index(path[1:], "/")
	if ci == -1 {
		ctlName = path[1:]
	} else {
		ctlName = path[1 : ci+1]
	}

	action = GetAction(path, this._urlSuffix)
	if len(action) != 0 {
		if strings.ToLower(action) == ctlName {
			action = "Index"
		}

		if ctx.Request.Method == "POST" {
			action += "_post"
		}

		if this._ctlMap != nil {
			if v := this._ctlMap[ctlName]; v != nil {
				CustomHandle(v(), ctx, action)
				return
			}
		}
	}

	http.Error(ctx.Response, "404 page not found",
		http.StatusNotFound)
}

// 用普通方式注册控制器，由生成器为每个请求生成一个控制器实例。
func (this *Route) NormalRegister(name string, cg ControllerGenerate) {
	if this._ctlMap == nil {
		this._ctlMap = make(map[string]ControllerGenerate)
	}
	this._ctlMap[name] = cg
}

// 用单例方式注册控制器，所有的请求都共享一个控制器.
func (this *Route) Register(name string, c Controller) {
	var cg ControllerGenerate = func() Controller { return c }
	this.NormalRegister(name, cg)
}

// Get Controller
func (this *Route) GetController(name string) Controller {
	if this._ctlMap != nil {
		if v, ok := this._ctlMap[name]; ok {
			return v()
		}
	}
	return nil
}
