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
	"errors"
	"fmt"
	"github.com/ixre/gof/crypto"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/util"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	defaultSessionMaxAge int64 = 3600 * 12
)

var (
	_storage           storage.Interface
	_defaultCookieName string
	_xsrfCookie        bool
	_factory           *sessionFactory = &sessionFactory{}
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
	if ns._storage == nil {
		panic(errors.New("session storage not set"))
	}
	ns._storage.Get(_factory.getStorageKey(ns._sessionId),
		&ns._data)
	return ns
}

func newSession(w http.ResponseWriter, s storage.Interface,
	cookieName string) *Session {
	id := _factory.newSessionId(s)
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
	s._storage.Get(_factory.getStorageKey(s._sessionId), &s._data)
	s.flushToClient()
}

// 销毁会话
func (s *Session) Destroy() {
	s._data = nil
	s._storage.Del(_factory.getStorageKey(s._sessionId))
	s.setMaxAge(-1e9)
	s.flushToClient()
}

// 保存会话
func (s *Session) Save() error {
	if s._data == nil {
		return nil
	}
	err := s._storage.SetExpire(_factory.getStorageKey(s._sessionId),
		&s._data, s._maxAge)
	if err == nil {
		s.flushToClient()
	}
	return err
}

// 设置过期秒数
func (s *Session) setMaxAge(seconds int64) {
	s._maxAge = seconds
}

// 获取会话的保存时间
func (s *Session) MaxAge() int64 {
	return s._maxAge
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

// 返回新的XSRF令牌
func (s *Session) NewXSRFToken() string {
	if !_xsrfCookie {
		return ""
	}
	unix := time.Now().UnixNano()
	str := fmt.Sprintf("%d-%d", unix, len(s._data))
	token := crypto.Md5([]byte(str))
	d := time.Duration(s._maxAge * 1e9)
	expires := time.Now().Local().Add(d)
	ck := &http.Cookie{
		Name:    "_xsrf_token",
		Value:   token,
		Path:    "/",
		Expires: expires,
	}
	http.SetCookie(s._rsp, ck)
	s.Set("_xsrf_token", token)
	s.Save()
	return token
}

// 检查XSRF令牌
func (s *Session) CheckXSRFToken(token string) bool {
	if _xsrfCookie {
		src := s.Get("_xsrf_token")
		if token != "" && src != nil {
			return src.(string) == token
		}
		return false
	}
	return true
}
func (s *Session) ResetXSRFToken() {
	if _xsrfCookie {
		s.Remove("_xsrf_token")
		s.Save()
	}
}

type sessionFactory struct{}

// 获取Session存储的键
func (s *sessionFactory) getStorageKey(sessionId string) string {
	return "gof:ss:" + sessionId
}

// create new session id
func (s *sessionFactory) newSessionId(sto storage.Interface) string {
	var key string
	for {
		key = strconv.Itoa(int(time.Now().Unix())) +
			util.RandString(8)
		key = crypto.Md5([]byte(key))[8:24]
		if !sto.Exists(s.getStorageKey(key)) {
			break
		}
	}
	return key
}

func (s *sessionFactory) parseSession(rsp http.ResponseWriter, r *http.Request,
	cookieName string, sto storage.Interface) *Session {
	if s == nil {
		log.Fatalln("session storage is nil")
	}
	if ck, err := r.Cookie(cookieName); err == nil {
		return getSession(rsp, sto, ck.Name, ck.Value)
	}
	return newSession(rsp, sto, cookieName)
}

// Initialize global session storage and name
func Initialize(sto storage.Interface, defaultName string, xsrfCookie bool) {
	_storage = sto
	if defaultName == "" {
		_defaultCookieName = "gof_sessionId"
	} else {
		_defaultCookieName = defaultName
	}
	_xsrfCookie = xsrfCookie
	// register session type for gob
	gob.Register(make(map[string]interface{}))
}

// Session adapter for http context
func Default(rsp http.ResponseWriter, r *http.Request) *Session {
	return _factory.parseSession(rsp, r, _defaultCookieName, _storage)
}
