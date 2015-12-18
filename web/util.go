/**
 * Copyright 2015 @ z3q.net.
 * name : util
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package web

import (
	"net/http"
	"strconv"
	"time"
)

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
