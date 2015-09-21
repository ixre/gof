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
	"net/http"
	"runtime"
	"runtime/debug"
)

func HttpError(rsp http.ResponseWriter, err error) {
	_, f, line, _ := runtime.Caller(1)
	rsp.Header().Add("Content-Type", "text/html")
	rsp.WriteHeader(500)

	var part1 string = `<html><head><title>HTTP ERROR </title>
				<meta charset="utf-8"/>
				<style>
				body{background:#FFF;font-size:100%;color:#333;margin:0 2%;}
        h1{color:red;font-size:28px;border-bottom:solid 1px #ddd;line-height:80px;}
        div.except-panel p{margin:20px 0;}
        div.except-panel div.summary{}
        div.except-panel p.message{font-size:24px;}
        div.except-panel p.contact{color:#666;font-size:18px;}
        div.except-panel p.stack{padding-top:30px;}
        div.except-panel p.stack em{font-size:18px;font-style: normal;}
        div.except-panel pre{font-family: Sans,Arail;
            border:solid 1px #ddd;padding:20px;
            font-size:16px;background:#F5F5F5;
            line-height: 150%;color:#888;}
        div.except-panel .hidden{display:none;}
			</style>
        </head>
        <body>`

	var html string = fmt.Sprintf(`
				<h1>HTTP ERROR ：%s</h1>
				<div class="except-panel">
					<div class="summary">
						<p class="message">Source：%s&nbsp;&nbsp;Line:%d</p>
						<p class="contact">Plese contact administrator or <a href="/">go home</a></p>
					</div>
					<p class="stack">
						<em>Stack：</em><br/>
						<pre>
							%s
						</pre>
					</p>
				</div>
		</body>
		</html>
		`, err.Error(), f, line, debug.Stack())

	rsp.Write([]byte(part1 + html))
}
