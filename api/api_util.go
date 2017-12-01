// create for mzl-api 07/11/2017 ( jarrysix@gmail.com )
package api

import (
	"github.com/jsix/gof"
	"github.com/jsix/gof/util"
	"reflect"
)

//转换到实体
func FormEntity(form Form, dst interface{}) (err error) {
	refVal := reflect.ValueOf(dst).Elem()
	for k, v := range form {
		d := refVal.FieldByName(k)
		if !d.IsValid() {
			continue
		}
		err = gof.AssignValue(d, util.Str(v))
	}
	return err
}
