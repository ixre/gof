/**
 * Copyright 2015 @ S1N1 Team.
 * name : context.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */

package web

import (
	"github.com/jrsix/gof"
	"net/http"
)

var globalSessionStorage gof.Storage

// 设置全局的会话存储
func SetSessionStorage(s gof.Storage) {
	globalSessionStorage = s
}

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
	if globalSessionStorage == nil {
		return this.App.Storage()
	}
	return globalSessionStorage
}

func (this *Context) Session() *Session {
	if this._session == nil {
		ck, err := this.Request.Cookie(sessionCookieName)
		ss := this.getSessionStorage()
		if err == nil {
			this._session = LoadSession(this.Response, ss, ck.Value)
		} else {
			this._session = NewSession(this.Response, ss)
			this._session.Save()
		}
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
