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
	"github.com/jsix/gof/storage"
	"github.com/jsix/gof/util"
	"log"
	"net/http"
	"time"
)

const (
	defaultSessionMaxAge int64 = 3600 * 12
)

var (
	_storage           storage.Interface
	_defaultCookieName string = "gof_SessionId"
)

type Session struct {
	_sessionId string
	_rsp       http.ResponseWriter
	_data      map[string]interface{}
	_storage   storage.Interface
	_maxAge    int64
	_keyName   string
}

func getSession(w http.ResponseWriter, s storage.Interface,
	cookieName string, sessionId string) *Session {
	ns := &Session{
		_sessionId: sessionId,
		_rsp:       w,
		_data:      make(map[string]interface{}),
		_storage:   s,
		_maxAge:    defaultSessionMaxAge,
		_keyName:   cookieName,
	}
	ns._storage.Get(getSessionId(ns._sessionId), &ns._data)
	return ns
}

func newSession(w http.ResponseWriter, s storage.Interface, cookieName string) *Session {
	id := newSessionId(s)
	return &Session{
		_sessionId: id,
		_rsp:       w,
		_storage:   s,
		_maxAge:    defaultSessionMaxAge,
		_keyName:   cookieName,
	}
}

func (s *Session) chkInit() {
	if s._data == nil {
		s._data = make(map[string]interface{})
	}
}

// 获取会话编号
func (s *Session) GetSessionId() string {
	return s._sessionId
}

// 获取值
func (s *Session) Get(key string) interface{} {
	if s._data != nil {
		if v, ok := s._data[key]; ok {
			return v
		}
	}
	return nil
}

// 设置键值
func (s *Session) Set(key string, v interface{}) {
	s.chkInit()
	//	if reflect.TypeOf(v).Kind() == reflect.Ptr{
	//		panic("Session value must be ptr")
	//	}
	s._data[key] = v
}

// 移除键
func (s *Session) Remove(key string) bool {
	if _, exists := s._data[key]; exists {
		delete(s._data, key)
		return true
	}
	return false
}

// 使用指定的会话代替当前会话
func (s *Session) UseInstead(sessionId string) {
	s._sessionId = sessionId
	s._storage.Get(getSessionId(s._sessionId), &s._data)
	s.flushToClient()
}

// 销毁会话
func (s *Session) Destroy() {
	s._storage.Del(getSessionId(s._sessionId))
	s.SetMaxAge(-s._maxAge)
	s.flushToClient()
}

// 保存会话
func (s *Session) Save() error {
	if s._data == nil {
		return nil
	}
	err := s._storage.SetExpire(getSessionId(s._sessionId), &s._data, s._maxAge)
	if err == nil {
		s.flushToClient()
	}
	return err
}

// 设置过期秒数
func (s *Session) SetMaxAge(seconds int64) {
	s._maxAge = seconds
}

//存储到客户端
func (s *Session) flushToClient() {
	d := time.Duration(s._maxAge * 1e9)
	expires := time.Now().Local().Add(d)
	ck := &http.Cookie{
		Name:     s._keyName,
		Value:    s._sessionId,
		Path:     "/",
		HttpOnly: true,
		Expires:  expires,
	}
	http.SetCookie(s._rsp, ck)
}

func init() {
	// register session type for gob.
	gob.Register(make(map[string]interface{}))
}

func getSessionId(id string) string {
	return "gof:ss:" + id
}

func newSessionId(s storage.Interface) string {
	var rdStr string
	for {
		dt := time.Now()
		randStr := util.RandString(6)
		rdStr = fmt.Sprintf("%s%d", randStr, dt.Second())
		if !s.Exists(getSessionId(rdStr)) {
			//check session id exists
			break
		}
	}
	return rdStr
}

// Set global session storage and name
func Set(s storage.Interface, defaultName string) {
	_storage = s
	if len(defaultName) > 0 {
		_defaultCookieName = defaultName
	}
}

func parseSession(rsp http.ResponseWriter, r *http.Request,
	cookieName string, s storage.Interface) *Session {
	if s == nil {
		log.Fatalln("session storage is nil")
	}
	if ck, err := r.Cookie(cookieName); err == nil {
		return getSession(rsp, s, ck.Name, ck.Value)
	}
	return newSession(rsp, s, cookieName)
}

// Session adapter for http context
func Default(rsp http.ResponseWriter, r *http.Request) *Session {
	return parseSession(rsp, r, _defaultCookieName, _storage)
}
