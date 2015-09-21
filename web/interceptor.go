package web

import (
	"errors"
	"fmt"
	"github.com/jsix/gof"
	"github.com/jsix/gof/log"
	"net/http"
	"os"
	"time"
)

var (
	HandleDefaultHttpExcept func(*Context, error)
	HandleHttpBeforePrint   func(*Context) bool
	HandleHttpAfterPrint    func(*Context)
	_                       http.Handler = new(Interceptor)
)

//Http请求处理代理
type Interceptor struct {
	_app gof.App
	//执行请求
	_execute RequestHandler
	//请求之前发生;返回false,则终止运行
	Before func(*Context) bool
	After  func(*Context)
	Except func(*Context, error)
}

func NewInterceptor(app gof.App, f RequestHandler) *Interceptor {
	return &Interceptor{
		_app:     app,
		_execute: f,
	}
}

func (this *Interceptor) handle(app gof.App, w http.ResponseWriter, r *http.Request, handler RequestHandler) {
	// proxy response writer
	//w := NewRespProxyWriter(w)
	ctx := NewContext(app, w, r)

	//todo: panic可以抛出任意对象，所以recover()返回一个interface{}
	if this.Except != nil {
		defer func() {
			if err := recover(); err != nil {
				this.Except(ctx, errors.New(fmt.Sprintf("%s", err)))
			}
		}()
	}

	if this.Before != nil {
		if !this.Before(ctx) {
			return
		}
	}
	if handler != nil {
		handler(ctx)
	}

	if this.After != nil {
		this.After(ctx)
	}
}

func (this *Interceptor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if this._app == nil {
		log.Fatalln("Please use web.NewInterceptor(gof.App) to initialize!")
		os.Exit(1)
	}
	this.handle(this._app, w, r, this._execute)
}

func (this *Interceptor) For(handle RequestHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		this.handle(this._app, w, r, handle)
	})
}

func init() {
	HandleDefaultHttpExcept = func(ctx *Context, err error) {
		HttpError(ctx.Response, err)
	}

	HandleHttpBeforePrint = func(ctx *Context) bool {
		r := ctx.Request
		fmt.Println("[Request] ", time.Now().Format("2006-01-02 15:04:05"), ": URL:", r.RequestURI)
		for k, v := range r.Header {
			fmt.Println(k, ":", v)
		}
		if r.Method == "POST" {
			r.ParseForm()
		}
		for k, v := range r.Form {
			fmt.Println("form", k, ":", v)
		}
		return true
	}

	HandleHttpAfterPrint = func(ctx *Context) {
		w := ctx.Response
		proxy, ok := w.ResponseWriter.(*ResponseProxyWriter)

		if !ok {
			fmt.Println("[Response] convert error")
			return
		}
		fmt.Println("[Respose]\n" + string(proxy.Output))
	}
}
