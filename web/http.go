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
