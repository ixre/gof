package gof

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

func AssignValue(d reflect.Value, s string) (err error) {
	switch d.Type().Kind() {
	case reflect.Float32, reflect.Float64:
		var x float64
		x, err = strconv.ParseFloat(s, d.Type().Bits())
		d.SetFloat(x)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var x int64
		x, err = strconv.ParseInt(s, 10, d.Type().Bits())
		d.SetInt(x)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var x uint64
		x, err = strconv.ParseUint(s, 10, d.Type().Bits())
		d.SetUint(x)
	case reflect.Bool:
		v := strings.ToLower(s)
		d.SetBool(v == "true" || v == "on" || v == "1")
	case reflect.String:
		d.SetString(s)
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
	return err

}
