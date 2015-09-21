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
	"net/http"
)

type Context struct {
	App      gof.App
	Request  *http.Request
	Response *response
	// 用于上下文数据交换
	Items    map[string]interface{}
	_session *Session
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
	}
}

func (this *Context) getSessionStorage() gof.Storage {
	if sessionStorage == nil {
		return this.App.Storage()
	}
	return sessionStorage
}

func (this *Context) Session() *Session {
	if this._session == nil {
		this._session = parseSession(this.Response, this.Request,
			sessionCookieName, this.getSessionStorage())
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
