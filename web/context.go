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
	Items    map[string]interface{} // 用于上下文数据交换
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
		_session: nil,
	}
}

func (c *Context) Session() *session.Session {
	if c._session == nil {
		c._session = session.Default(c.Response, c.Request)
	}
	return c._session
}

// 获取数据项
func (c *Context) GetItem(key string) interface{} {
	if v, e := c.Items[key]; e {
		return v
	}
	return nil
}
