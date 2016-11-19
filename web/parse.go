package web

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

//转换到实体
func ParseFormToEntity(values map[string][]string, instance interface{}) {
	refVal := reflect.ValueOf(instance).Elem()
	//类型装换参见：http://www.kankanews.com/ICkengine/archives/19245.shtml
	//for i:=0 ; i< refVal.NumField(); i++ {
	//	prop := refVal.Field(i)
	for k, v := range values {
		field := refVal.FieldByName(k)
		if field.IsValid() {
			//
			//var x interface{} = 1
			//y:= x.(type)
			//
			strVal := v[0]

			switch field.Type().Kind() {
			case reflect.String:
				field.Set(reflect.ValueOf(strVal))

			case reflect.Float32:
				if v, err := strconv.ParseFloat(strVal, 32); err == nil && v != 0 {
					field.Set(reflect.ValueOf(float32(v)))
				}
			case reflect.Float64:
				if v, err := strconv.ParseFloat(strVal, 64); err == nil && v != 0 {
					field.Set(reflect.ValueOf(v))
				}
			case reflect.Int8:
				val, err := strconv.ParseUint(strVal, 10, 8)
				if err == nil {
					field.Set(reflect.ValueOf(int8(val)))
				}
			case reflect.Int16:
				val, err := strconv.ParseUint(strVal, 10, 16)
				if err == nil {
					field.Set(reflect.ValueOf(int16(val)))
				}
			case reflect.Int32:
				val, err := strconv.Atoi(strVal)
				if err == nil {
					field.Set(reflect.ValueOf(int32(val)))
				}
			case reflect.Int64:
				val, err := strconv.Atoi(strVal)
				if err == nil {
					field.Set(reflect.ValueOf(int64(val)))
				}
			case reflect.Int:
				val, err := strconv.Atoi(strVal)
				if err == nil {
					field.Set(reflect.ValueOf(val))
				}
			case reflect.Bool:
				val := strings.ToLower(strVal) == "true" || strVal == "1" || strVal == "on"
				field.Set(reflect.ValueOf(val))

			case reflect.Struct:
				v := field.Interface()
				switch v.(type) {
				case time.Time:
					t, err := time.Parse("2006-01-02 15:04:05", strVal)
					if err == nil {
						field.Set(reflect.ValueOf(t))
					}
				}
			default:
				panic("not support type:" + field.Type().Kind().String())
			}

			//接口类型
			//			case reflect.Interface:
			//				if reflect.TypeOf(time.Now()) == field.Type() {
			//					t, err := time.Parse("2006-01-02 15:04:05", strVal)
			//					if err == nil {
			//						field.Set(reflect.ValueOf(t))
			//					}
			//				}
			//				break
			//			}
		}
	}
	//fmt.Println(instance)
}
