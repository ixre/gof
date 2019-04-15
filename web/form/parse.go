package form

import (
	"github.com/ixre/gof/net/http"
)

//转换到实体
func ParseEntity(values map[string][]string, dst interface{}) (err error) {
	return http.MapToEntity(values, dst)
}
