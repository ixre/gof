/**
 * Copyright 2015 @ z3q.net.
 * name : context.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */

package web

import (
	"github.com/jsix/gof"
	"github.com/jsix/gof/web/session"
	"net/http"
)

type Context struct {
	App      gof.App
	Request  *http.Request
	Response *response
	// 用于上下文数据交换
	Items    map[string]interface{}
	_session *session.Session
}

func NewContext(app gof.App, rsp http.ResponseWriter, req *http.Request) *Context {
	newRsp := &response{
		ResponseWriter: rsp,
	}
	return &Context{
		App:      app,
		Response: newRsp,
		Request:  req,
		Items:    make(map[string]interface{}),
		_session: app.Storage(),
	}
}

func (this *Context) Session() *session.Session {
	if this._session == nil {
		this._session = session.Default(this.Response, this.Request)
	}
	return this._session
}

// 获取数据项
func (this *Context) GetItem(key string) interface{} {
	if v, e := this.Items[key]; e {
		return v
	}
	return nil
}
