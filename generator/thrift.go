/**
 * Copyright 2015 @ at3.net.
 * name : thrift.go
 * author : jarryliu
 * date : 2016-11-17 13:14
 * description :
 * history :
 */
package generator

import (
	"bytes"
	"errors"
	"reflect"
	"strconv"
)

// 转换结构为Thrift的结构
func ThriftStruct(v interface{}) ([]byte, error) {
	vt := reflect.TypeOf(v)
	if vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
	}
	if vt.Kind() != reflect.Struct {
		return nil, errors.New("value is not struct")
	}
	buf := bytes.NewBufferString("")
	buf.WriteString("struct ")
	buf.WriteString(vt.Name())
	buf.WriteString(" {\n")
	for i, n := 0, vt.NumField(); i < n; i++ {
		f := vt.Field(i)
		buf.WriteString("    ")
		buf.WriteString(strconv.Itoa(i + 1))
		buf.WriteString(": ")
		buf.WriteString(thriftType(f.Type))
		buf.WriteString(" ")
		buf.WriteString(f.Name)
		buf.WriteString("\n")
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func thriftType(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Int8, reflect.Uint8:
		return "i8"
	case reflect.Int16, reflect.Uint16:
		return "i16"
	case reflect.Int32, reflect.Uint32:
		return "i32"
	case reflect.Int, reflect.Uint, reflect.Int64, reflect.Uint64:
		return "i64"
	case reflect.Array:
		return "slist"
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "bool"
	case reflect.Float32, reflect.Float64:
		return "double"
	}
	return "binary"
}
