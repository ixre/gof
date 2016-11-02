/**
 * Copyright 2015 @ z3q.net.
 * name : error
 * author : jarryliu
 * date : 2015-09-21 11:22
 * description :
 * history :
 */
package web

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

type (
	// 一个针对多个子域的HTTP处理程序
	MultiHttpHandler interface {
		// 设置默认的处理程序
		Default(handler http.Handler)
		// 添加子域的处理程序
		Set(sub string, handler http.Handler)
		// 获取处理程序
		Get(sub string) http.Handler
		// 处理HTTP请求
		ServeHTTP(w http.ResponseWriter, r *http.Request)
		// 监听端口,并启动
		ListenAndServe(addr string) error
	}

	HttpHostsHandler map[string]http.Handler
)

var _ MultiHttpHandler = new(HttpHostsHandler)

func (h HttpHostsHandler) ListenAndServe(addr string) error {
	log.Println("** server running on", addr)
	err := http.ListenAndServe(addr, h)
	if err != nil {
		log.Println("** serve exit! ", err.Error())
	}
	return err
}

func (h HttpHostsHandler) Default(handler http.Handler) {
	h["*"] = handler
}

// 获取处理程序
func (h HttpHostsHandler) Get(subName string) http.Handler {
	return h[subName]
}

func (h HttpHostsHandler) Set(subName string, handler http.Handler) {
	h[subName] = handler
}

func (h HttpHostsHandler) GetSubName(r *http.Request) string {
	return r.Host[:strings.Index(r.Host, ".")+1]
}

func (h HttpHostsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hh, ok := h[h.GetSubName(r)] //根据主机头返回响应内容
	if !ok {
		hh, _ = h["*"] //获取通用的serve
	}
	if hh != nil {
		hh.ServeHTTP(w, r)
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func HttpError(rsp http.ResponseWriter, err error) {
	_, f, line, _ := runtime.Caller(1)
	rsp.Header().Add("Content-Type", "text/html")
	rsp.WriteHeader(500)

	var part1 string = `<html><head><title>HTTP ERROR</title>
				<meta charset="utf-8"/>
				<style>
	body{background:#FFF;font-size:100%;margin:0 0 2em 0;}
	div.tit{background:#FFA;border-bottom:solid 1px #FE0;}
    h1{color:#F00;font-size:2em;line-height:2em;margin:0;border-bottom:solid 1px #FFF;padding:0.8em 2%;}
    div.except-panel p{margin:0 2%;padding:0}
    div.except-panel div.summary{color:#000;}
    div.except-panel p.message{font-size:1.4em;margin:2em 2% 0 2%;}
    div.except-panel p.contact{color:#666;font-size:1.2em;}
    div.except-panel p.stack{padding:0;}
    div.except-panel em{font-size:1.2em;font-style: normal;color:#000;line-height:2em;font-weight:bold;}
    div.except-panel pre{font-family: Sans,Arail;margin:1em 2% 2em 2%;line-height: 150%;color:#333;}
    div.except-panel .hidden{display:none;}
			</style>
        </head>
        <body>`

	var html string = fmt.Sprintf(`
				<div class="tit"><h1>%s</h1></div>
				<div class="except-panel">
					<div class="summary">
						<p class="message">Source：%s&nbsp;&nbsp;Line:%d</p>
					</div>
					<p class="stack">
						<pre><em>Stack：</em><br/>%s</pre>
					</p>
					<p class="contact">
						<p class="contact">Plese contact administrator or [ <a href="/">Go Home</a> ]</p>
					</p>
				</div>
		</body>
		</html>
		`, err.Error(), f, line, debug.Stack())

	rsp.Write([]byte(part1 + html))
}

// 设置缓存头部信息
func SetCacheHeader(w http.ResponseWriter, minute int) {
	h := w.Header()
	t := time.Now()
	expires := time.Minute * time.Duration(minute)
	h.Set("Pragma", "Pragma")                 //Pragma:设置页面是否缓存，为Pragma则缓存，no-cache则不缓存
	h.Set("Expires", t.Add(expires).String()) //Expires:过时期限值
	//h.Set("Last-Modified",t.String()); 			//Last-Modified:页面的最后生成时间
	h.Set("Cache-Control", "max-age="+strconv.Itoa(minute*60)) //Cache-Control来控制页面的缓存与否,public:浏览器和缓存服务器都可以缓存页面信息；
}
