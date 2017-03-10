package form

import (
	"github.com/jsix/gof/web"
)

//转换到实体
func ParseEntity(values map[string][]string, dst interface{}) (err error) {
	return web.ParseEntity(values, dst)
}
