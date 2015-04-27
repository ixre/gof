/**
 * Copyright 2015 @ S1N1 Team.
 * name : route.go
 * author : newmin
 * date : -- :
 * description :
 * history :
 */
package mvc

import (
    "github.com/atnet/gof/web"
    "strings"
    "net/http"
)

var _ web.Route = new(Route)

type Route struct{
    _lazyRegisted bool
    _routeMap *web.RouteMap
    _ctlMap map[string]interface{}
}

func NewRoute(source *web.RouteMap)*Route{
    if source == nil {
        source = new(web.RouteMap)
    }

    r := &Route{
        _ctlMap : make(map[string]interface{}),
        _routeMap : source,
    }
    return r
}

//添加路由
func (this *Route) Add(urlPattern string, rf web.HttpContextFunc) {
    this._routeMap.Add(urlPattern,this.defaultRouteHandle)
}

// 处理请求
func (this *Route) Handle(ctx *web.Context) {
    if !this._lazyRegisted {
        // 添加默认的路由
        this.Add("*",this.defaultRouteHandle)
        this._lazyRegisted = true
    }
    this._routeMap.Handle(ctx)
}

func (this *Route) defaultRouteHandle(ctx *web.Context) {
    path := ctx.Request.URL.Path
    var ctlName, action string
    ci := strings.Index(path[1:], "/")
    if ci == -1 {
        ctlName = path[1:]
        action = "Index"
    }else {
        ctlName = path[1:ci]
        path = path[ci+1:]
        ai := strings.Index(path, "/")
        if ai == -1 {
            action = path
        }else {
            action= path[:ai]
        }
    }

    if this._ctlMap != nil {
        if v := this._ctlMap[ctlName]; v != nil {
            CustomHandle(v, ctx, strings.ToTitle(action), nil)
            return
        }
    }

    http.Error(ctx.ResponseWriter, "404 page not found",
    http.StatusNotFound)
}

// 注册控制器
func (this *Route) RegisterController(name string,ctl interface{}){
    if this._ctlMap == nil{
        this._ctlMap = make(map[string]interface{})
    }
    this._ctlMap[name] = ctl
}