/**
 * Copyright 2015 @ S1N1 Team.
 * name : route.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package mvc

import (
	"github.com/atnet/gof/web"
	"log"
	"net/http"
	"strings"
)

var _ web.Route = new(Route)

type Route struct {
	_lazyRegister bool
	_routeMap     *web.RouteMap
	_ctlMap       map[string]Controller
}

func NewRoute(source *web.RouteMap) *Route {
	if source == nil {
		source = new(web.RouteMap)
	}

	r := &Route{
		_ctlMap:   make(map[string]Controller),
		_routeMap: source,
	}
	return r
}

//添加路由
func (this *Route) Add(urlPattern string, rf web.HttpContextFunc) {
	if urlPattern == "" || urlPattern == "/" || urlPattern == "^/" {
		log.Fatalln("[ Panic] - Dangerous!The url parttern \"" + urlPattern +
			"\" will override default route," +
			"please instead of \"^/$\"")
	}
	this._routeMap.Add(urlPattern, rf)
}

// 处理请求
func (this *Route) Handle(ctx *web.Context) {
	if !this._lazyRegister {
		// 添加默认的路由
		this.Add("*", this.defaultRouteHandle)
		this._lazyRegister = true
	}
	this._routeMap.Handle(ctx)
}

func (this *Route) defaultRouteHandle(ctx *web.Context) {
	path := ctx.Request.URL.Path
	var ctlName, action string
	ci := strings.Index(path[1:], "/")
	if ci == -1 {
		ctlName = path[1:]
	} else {
		ctlName = path[1 : ci+1]
		path = path[ci+2:]
		ai := strings.Index(path, "/")
		if ai == -1 {
			action = path
		} else {
			action = path[:ai]
		}
	}
	if len(action) == 0 {
		action = "Index"
	} else {
		//将第一个字符转为大写,这样才可以
		upperFirstLetter := strings.ToUpper(action[0:1])
		if upperFirstLetter != action[0:1] {
			action = upperFirstLetter + action[1:]
		}
	}

	if ctx.Request.Method == "POST" {
		action = action + "_post"
	}

	if this._ctlMap != nil {
		if v := this._ctlMap[ctlName]; v != nil {
			CustomHandle(v, ctx, action)
			return
		}
	}

	http.Error(ctx.ResponseWriter, "404 page not found",
		http.StatusNotFound)
}

// Register Controller into routes table.
func (this *Route) RegisterController(name string, ctl Controller) {
	if this._ctlMap == nil {
		this._ctlMap = make(map[string]Controller)
	}
	this._ctlMap[name] = ctl
}

// Get Controller
func (this *Route) GetController(name string) Controller {
	if this._ctlMap != nil {
		if v, ok := this._ctlMap[name]; ok {
			return v
		}
	}
	return nil
}
