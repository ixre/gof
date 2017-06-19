package web

import (
	"github.com/jsix/gof"
	"github.com/jsix/gof/storage"
	"github.com/jsix/gof/web/session"
	"net/http"
	"reflect"
	"strings"
)

type Options struct {
	Storage           storage.Interface
	SessionCookieName string
	XSRFCookie        bool
}

func Initialize(o Options) {
	session.Initialize(o.Storage, o.SessionCookieName, o.XSRFCookie)
}

// 获取请求完整的地址
func RequestRawURI(r *http.Request) string {
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	return strings.Join([]string{scheme, r.Host, r.RequestURI}, "")
}

// 获取协议
func Scheme(r *http.Request) string {
	if r.TLS != nil {
		return "https://"
	}
	return "http://"
}

//转换到实体
func ParseEntity(values map[string][]string, dst interface{}) (err error) {
	refVal := reflect.ValueOf(dst).Elem()
	for k, v := range values {
		d := refVal.FieldByName(k)
		if !d.IsValid() {
			continue
		}
		err = gof.AssignValue(d, v[0])
	}
	return err
}
