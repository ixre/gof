/**
 * Copyright 2015 @ S1N1 Team.
 * name : session.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package web
import (
    "github.com/atnet/gof"
    "time"
    "fmt"
    "github.com/atnet/gof/util"
    "net/http"
)

const defaultSessionTime int = 7200
const  sessionCookieName string = "_gofs"

func getSessionKey(key string)string{
    return "gof:web:session:"+key
}

func newSessionKey()string{
    dt := time.Now()
    randStr := util.RandString(4)
    return fmt.Sprintf("%s%d%d",randStr,dt.Second(),dt.Nanosecond())
}

type Session struct {
    _key    string
    _rsp    http.ResponseWriter
    _data map[string]interface{}
    _storage gof.Storage
    _exipresSeconds int
}

func LoadSession(w http.ResponseWriter,storage gof.Storage,key string)*Session {
    s := &Session{
        _key : key,
        _rsp :w,
        _storage : storage,
        _exipresSeconds : defaultSessionTime,
    }
    s._storage.Get(getSessionKey(s._key), s._data)
    return s
}

func NewSession(w http.ResponseWriter,storage gof.Storage)*Session{
    key := newSessionKey()
    return &Session{
        _key : key,
        _rsp :w,
        _storage : storage,
        _exipresSeconds : defaultSessionTime,
    }
}

func (this *Session) chkInit(){
    if this._data == nil{
        this._data = make(map[string]interface{})
    }
}

func (this *Session) Get(key string)interface{}{
    if this._data!= nil{
        if v,ok := this._data[key];ok{
            return v
        }
    }
    return nil
}

func (this *Session) Set(key string,v interface{}){
    this.chkInit()
    this._data[key] = v
}

// 保存会话
func (this *Session) Save()error {
    if this._data == nil {
        return nil
    }
    err := this._storage.Set(getSessionKey(this._key), this._data)
    if err == nil{
        this.flushToClient()
    }
    return err
}

// 设置过期秒数
func (this *Session) SetExpires(seconds int){
    this._exipresSeconds = seconds
}

//存储到客户端
func (this *Session) flushToClient(){
    d := time.Duration(this._exipresSeconds)
    ck := &http.Cookie{
        Name: sessionCookieName,
        Value : this._key,
        Path:"/",
        Expires: time.Now().Add(d),
    }
    http.SetCookie(this._rsp,ck)
}