/**
 * Copyright 2015 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2015-04-27 00:53
 * description :
 * history :
 */
package main

import (
	"fmt"
	"github.com/atnet/gof"
	"github.com/atnet/gof/db"
	"github.com/atnet/gof/log"
	"github.com/atnet/gof/web"
	"github.com/atnet/gof/web/mvc"
	"net/http"
	"strings"
)

// Implement gof.App
type HttpApp struct {
	config      *gof.Config
	dbConnector db.Connector
	template    *gof.Template
	logger      log.ILogger
}

// init settings
func (this *HttpApp) Config() *gof.Config {
	if this.config == nil {
		this.config = gof.NewConfig()
		this.config.Set("SYS_NAME", "DEMO")
		this.config.Set("MYSQL_HOST", "127.0.0.1")
		this.config.Set("MYSQL_PORT", 3306)
		this.config.Set("MYSQL_MAXCONN", 1000)
		this.config.Set("MYSQL_USR", "root")
		this.config.Set("MYSQL_PWD", "")
		this.config.Set("MYSQL_DBNAME", "")
	}
	return this.config
}

// init database connector
func (this *HttpApp) Db() db.Connector {
	if this.dbConnector == nil {
		source := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf-8&loc=Local",
			this.Config().GetString("MYSQL_USR"),
			this.Config().GetString("MYSQL_PWD"),
			this.Config().GetString("MYSQL_HOST"),
			this.Config().GetInt("MYSQL_PORT"),
			this.Config().GetString("MYSQL_DBNAME"),
		)
		this.dbConnector = db.NewBasicConnector("mysql", source, this.Log(),
			this.Config().GetInt("MYSQL_MAXCONN"))
	}
	return this.dbConnector
}

// init template
func (this *HttpApp) Template() *gof.Template {
	if this.template == nil {
		this.template = &gof.Template{
			Init: func(v *map[string]interface{}) {
				(*v)["SYS_NAME"] = this.Config().GetString("SYS_NAME")
			},
		}
	}
	return this.template
}

func (this *HttpApp) Source() interface{} {
	return this
}

func (this *HttpApp) Debug() bool {
	return false
}

func (this *HttpApp) Log() log.ILogger {
	if this.logger == nil {
		var flag int = 0
		if this.Debug() {
			flag = log.LOpen | log.LESource | log.LStdFlags
		}
		this.logger = log.NewLogger(nil, " O2O", flag)
	}
	return this.logger
}

//获取Http请求代理处理程序
func getInterceptor(a gof.App, routes web.Route) *web.Interceptor {
	var in = web.NewInterceptor(a, getHttpExecFunc(routes))
	in.Except = web.HandleDefaultHttpExcept
	in.Before = nil
	in.After = nil
	return in
}

// 获取执行方法
func getHttpExecFunc(routes web.Route) web.ContextFunc {
	return func(ctx *web.Context) {
		r, w := ctx.Request, ctx.ResponseWriter
		switch host, f := r.Host, strings.HasPrefix; {
		//静态文件
		case f(host, "static."):
			http.ServeFile(w, r, "./static"+r.URL.Path)

		default:
			routes.Handle(ctx)
		}
	}
}

// Test Controller
type testController struct {
}

func (this *testController) Requesting(ctx *web.Context) bool {
	ctx.ResponseWriter.Write([]byte("\r\nit's pass by filter...."))
	return !false
}

func (this *testController) RequestEnd(ctx *web.Context) {
	ctx.ResponseWriter.Write([]byte("\r\nRequest End."))
}

func (this *testController) Index(ctx *web.Context) {
	ctx.ResponseWriter.Write([]byte("\r\nRequesting....."))
}

func main() {
	app := &HttpApp{}
	routes := mvc.NewRoute(nil)
	routes.RegisterController("test", &testController{})

	routes.Add("/[0-9]$", func(ctx *web.Context) {
		ctx.ResponseWriter.Write([]byte("数字路径"))
	})

	routes.Add("/[a-z]$", func(ctx *web.Context) {
		ctx.ResponseWriter.Write([]byte("字母路径"))
	})

	routes.Add("^/$", func(ctx *web.Context) {
		sysName := ctx.App.Config().GetString("SYS_NAME")
		ctx.ResponseWriter.Write([]byte("Hello,Gof with " + sysName + "."))
		ctx.ResponseWriter.Header().Set("Content-Type", "text/html")
		return
		ctx.App.Template().ExecuteIncludeErr(ctx.ResponseWriter,
			func(v *map[string]interface{}) {
				(*v)["变量名"] = "变量值"
			}, "views/index.html")
	})

	var in = getInterceptor(app, routes)
	go http.ListenAndServe(":8080", in)

	log.Println("[ OK] - web is listening on port :8080.")
	var ch = make(chan int, 1)
	<-ch
}
