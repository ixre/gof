package form

import (
	"reflect"
	"strconv"
	"time"
)

//转换到实体
func ParseEntity(values map[string][]string, dst interface{}) (err error) {
	refVal := reflect.ValueOf(dst).Elem()
	//类型装换参见：http://www.kankanews.com/ICkengine/archives/19245.shtml
	//for i:=0 ; i< refVal.NumField(); i++ {
	//	prop := refVal.Field(i)
	for k, v := range values {
		d := refVal.FieldByName(k)
		if !d.IsValid() {
			continue
		}
		s := v[0]
		switch d.Type().Kind() {
		case reflect.Float32, reflect.Float64:
			var x float64
			x, err = strconv.ParseFloat(string(s), d.Type().Bits())
			d.SetFloat(x)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var x int64
			x, err = strconv.ParseInt(string(s), 10, d.Type().Bits())
			d.SetInt(x)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			var x uint64
			x, err = strconv.ParseUint(string(s), 10, d.Type().Bits())
			d.SetUint(x)
		case reflect.Bool:
			var x bool
			x, err = strconv.ParseBool(string(s))
			d.SetBool(x)
		case reflect.String:
			d.SetString(string(s))
		case reflect.Struct:
			v := d.Interface()
			switch v.(type) {
			case time.Time:
				t, err := time.Parse("2006-01-02 15:04:05", s)
				if err == nil {
					d.Set(reflect.ValueOf(t))
				}
			}
		}
		if err != nil {
			return err
		}
	}
	return err
}
