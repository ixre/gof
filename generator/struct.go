/**
 * Copyright 2015 @ at3.net.
 * name : struct
 * author : jarryliu
 * date : 2016-11-17 13:44
 * description :
 * history :
 */
package generator

import (
	"bytes"
	"errors"
	"reflect"
)

// 生成结构赋值代码
func StructAssignCode(v interface{}) ([]byte, error) {
	vt := reflect.TypeOf(v)
	if vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
	}
	if vt.Kind() != reflect.Struct {
		return nil, errors.New("value is not struct")
	}
	buf := bytes.NewBufferString("v := &")
	buf.WriteString(vt.Name())
	buf.WriteString(" {\n")
	for i, n := 0, vt.NumField(); i < n; i++ {
		f := vt.Field(i)
		buf.WriteString("    ")
		buf.WriteString(f.Name)
		buf.WriteString(" : src.")
		buf.WriteString(f.Name)
		buf.WriteString(",\n")
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}
