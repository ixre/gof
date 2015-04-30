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
	"github.com/atnet/gof"
	"net/http"
)

type HttpContextFunc func(*Context)

var sessionStorage gof.Storage

// 设置全局的会话存储
func SetSessionStorage(s gof.Storage){
	sessionStorage = s
}

type Context struct {
	App            gof.App
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	_session	   *Session
}

func NewContext(app gof.App, rsp http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		App:            app,
		ResponseWriter: rsp,
		Request:        req,
	}
}

func (this *Context) getSessionStorage()gof.Storage{
	if sessionStorage == nil{
		return this.App.Storage()
	}
	return sessionStorage
}

func (this *Context) Session()*Session{
	if this._session == nil {
		ck, err := this.Request.Cookie(sessionCookieName)
		ss := this.getSessionStorage()
		if err == nil {
			this._session = LoadSession(this.ResponseWriter, ss, ck.Value)
		}else {
			this._session = NewSession(this.ResponseWriter, ss)
			this._session.Save()
		}
	}
	return this._session
}