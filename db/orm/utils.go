package orm

import "reflect"

// 根据反射类型,返回orm对应的类型
func GetReflectTypeId(t reflect.Type) int {
	return int(t.Kind())
}
