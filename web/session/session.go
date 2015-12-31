/**
 * Copyright 2015 @ z3q.net.
 * name : session.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package session

import (
	"encoding/gob"
	"fmt"
	"github.com/jsix/gof"
	"github.com/jsix/gof/util"
	"log"
	"net/http"
	"time"
)

const (
	defaultSessionMaxAge int64 = 3600 * 12
)

var (
	_storage           gof.Storage
	_defaultCookieName string = "gof_SessionId"
)

type Session struct {
	_sessionId string
	_rsp       http.ResponseWriter
	_data      map[string]interface{}
	_storage   gof.Storage
	_maxAge    int64
	_keyName   string
}

func getSession(w http.ResponseWriter, storage gof.Storage, cookieName string, sessionId string) *Session {
	s := &Session{
		_sessionId: sessionId,
		_rsp:       w,
		_data:      make(map[string]interface{}),
		_storage:   storage,
		_maxAge:    defaultSessionMaxAge,
		_keyName:   cookieName,
	}
	s._storage.Get(getSessionId(s._sessionId), &s._data)
	return s
}

func newSession(w http.ResponseWriter, storage gof.Storage, cookieName string) *Session {
	id := newSessionId()
	return &Session{
		_sessionId: id,
		_rsp:       w,
		_storage:   storage,
		_maxAge:    defaultSessionMaxAge,
		_keyName:   cookieName,
	}
}

func (this *Session) chkInit() {
	if this._data == nil {
		this._data = make(map[string]interface{})
	}
}

// 获取会话编号
func (this *Session) GetSessionId() string {
	return this._sessionId
}

// 获取值
func (this *Session) Get(key string) interface{} {
	if this._data != nil {
		if v, ok := this._data[key]; ok {
			return v
		}
	}
	return nil
}

// 设置键值
func (this *Session) Set(key string, v interface{}) {
	this.chkInit()
//	if reflect.TypeOf(v).Kind() == reflect.Ptr{
//		panic("Session value must be ptr")
//	}
	this._data[key] = v
}

// 移除键
func (this *Session) Remove(key string) bool {
	if _, exists := this._data[key]; exists {
		delete(this._data, key)
		return true
	}
	return false
}

// 使用指定的会话代替当前会话
func (this *Session) UseInstead(sessionId string) {
	this._sessionId = sessionId
	this._storage.Get(getSessionId(this._sessionId), &this._data)
	this.flushToClient()
}

// 销毁会话
func (this *Session) Destroy() {
	this._storage.Del(getSessionId(this._sessionId))
	this.SetMaxAge(-this._maxAge)
	this.flushToClient()
}

// 保存会话
func (this *Session) Save() error {
	if this._data == nil {
		return nil
	}
	err := this._storage.SetExpire(getSessionId(this._sessionId), &this._data, this._maxAge)
	if err == nil {
		this.flushToClient()
	}
	return err
}

// 设置过期秒数
func (this *Session) SetMaxAge(seconds int64) {
	this._maxAge = seconds
}

//存储到客户端
func (this *Session) flushToClient() {
	d := time.Duration(this._maxAge * 1e9)
	expires := time.Now().Local().Add(d)
	ck := &http.Cookie{
		Name:     this._keyName,
		Value:    this._sessionId,
		Path:     "/",
		HttpOnly: true,
		Expires:  expires,
	}
	http.SetCookie(this._rsp, ck)
}

func init() {
	gob.Register(make(map[string]interface{})) // register session type for gob.
}

func getSessionId(id string) string {
	return "gof:session:" + id
}

func newSessionId() string {
	dt := time.Now()
	randStr := util.RandString(4)
	return fmt.Sprintf("%s%d%d", randStr, dt.Second(), dt.Nanosecond())
}

// Set global session storage and name
func Set(s gof.Storage, defaultName string) {
	_storage = s
	if len(defaultName) > 0 {
		_defaultCookieName = defaultName
	}
}

// get session storage
func getStorage() gof.Storage {
	return _storage
}

func parseSession(rsp http.ResponseWriter, r *http.Request,
	cookieName string, sto gof.Storage) *Session {
	if sto == nil {
		log.Fatalln("session storage is nil")
	}
	if ck, err := r.Cookie(cookieName); err == nil {
		return getSession(rsp, sto, ck.Name, ck.Value)
	}
	return newSession(rsp, sto, cookieName)
}

// Session adapter for http context
func Default(rsp http.ResponseWriter, r *http.Request) *Session {
	return parseSession(rsp, r, _defaultCookieName, _storage)
}

// create a session use custom key
func Create(key string, rsp http.ResponseWriter, r *http.Request) *Session {
	return parseSession(rsp, r, key, _storage)
}
