package http

import (
	"github.com/ixre/gof"
	"net"
	"net/http"
	"reflect"
	"strings"
)

// MapToEntity 转换到实体
func MapToEntity(values map[string][]string, dst interface{}) (err error) {
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

// RealIp 获取HTTP请求真实IP
func RealIp(r *http.Request) string {
	ra := r.RemoteAddr
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		ra = strings.Split(ip, ", ")[0]
	} else if ip := r.Header.Get("X-Real-IP"); ip != "" {
		ra = ip
	} else {
		ra, _, _ = net.SplitHostPort(ra)
	}
	return ra
}
