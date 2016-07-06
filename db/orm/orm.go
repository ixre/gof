package orm

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

type (
	TableMeta struct {
		TableName   string
		PkFieldName string
		PkIsAuto    bool
		//the index of fields
		FieldsIndex   []int
		FieldMapNames []string
	}

	Orm interface {
		// version of orm
		Version() string

		//Set orm output information
		SetTrace(b bool)

		// create the mapping data table
		Mapping(v interface{}, table string)

		Get(primaryVal interface{}, dst interface{}) error

		//get entity by condition
		GetBy(dst interface{}, where string, args ...interface{}) error

		//get entity by sql query result
		GetByQuery(dst interface{}, sql string, args ...interface{}) error

		//Select more than 1 entity list
		//@to : refrence to queryed entity list
		//@params : query condition
		//@where : other condition
		Select(dst interface{}, where string, args ...interface{}) error

		//select by sql query,dst must be one slice.
		SelectByQuery(dst interface{}, sql string, args ...interface{}) error

		//delete entity and effect to database
		Delete(entity interface{}, where string, args ...interface{}) (effect int64, err error)

		//delete entity by primary key
		DeleteByPk(entity interface{}, primary interface{}) (err error)

		Save(primary interface{}, entity interface{}) (rows int64, lastInsertId int64, err error)
	}
)

// 获取表元数据
func GetTableMapMeta(t reflect.Type) *TableMeta {
	ixs, maps := GetFields(t)
	pkName, pkIsAuto := GetPKName(t)
	m := &TableMeta{
		TableName:     t.Name(),
		PkFieldName:   pkName,
		PkIsAuto:      pkIsAuto,
		FieldsIndex:   ixs,
		FieldMapNames: maps,
	}
	return m
}

//if not defined primary key.the first key will as primary key
func GetPKName(t reflect.Type) (pkName string, pkIsAuto bool) {
	var ti int = t.NumField()

	ffc := func(f reflect.StructField) (string, bool) {
		if f.Tag != "" {
			var isAuto bool
			var fieldName string

			if ia := f.Tag.Get("auto"); ia == "yes" || ia == "1" {
				isAuto = true
			}

			if fieldName = f.Tag.Get("db"); fieldName != "" {
				return fieldName, isAuto
			}
			return f.Name, isAuto
		}
		return f.Name, false
	}

	for i := 0; i < ti; i++ {
		f := t.Field(i)
		if f.Tag != "" {
			pk := f.Tag.Get("pk")
			if pk == "1" || pk == "yes" {
				return ffc(f)
			}
		}
	}

	return ffc(t.Field(0))
}

// 获取实体的字段
func GetFields(t reflect.Type) (posArr []int, mapNames []string) {
	posArr = []int{}
	mapNames = []string{}

	fNum := t.NumField()
	var fmn string

	for i := 0; i < fNum; i++ {
		f := t.Field(i)
		if f.Tag != "" {
			fmn = f.Tag.Get("db")
			if fmn == "-" || fmn == "_" || len(fmn) == 0 {
				continue
			}
		}
		if fmn == "" {
			fmn = f.Name
		}
		internalKeysCheck(&fmn)
		mapNames = append(mapNames, fmn)
		posArr = append(posArr, i)
		fmn = ""
	}
	return posArr, mapNames
}

// format internal keywords
func internalKeysCheck(field *string) {
	switch *field {
	case "key", "where", "type":
		*field = strings.Join([]string{"`", *field, "`"}, "")
	}
}

func BindFields(meta *TableMeta, dst *reflect.Value, rawBytes [][]byte) error {
	for i, fi := range meta.FieldsIndex {
		SetField(dst.Field(fi), rawBytes[i])
	}
	return nil
}

func SetField(field reflect.Value, d []byte) {
	if field.IsValid() {
		//fmt.Println(field.String(), "==>", field.Type().Kind())
		switch field.Type().Kind() {
		case reflect.String:
			field.Set(reflect.ValueOf(string(d)))
			return

		case reflect.Int:
			val, err := strconv.ParseInt(string(d), 10, 0)
			if err == nil {
				field.Set(reflect.ValueOf(int(val)))
			}
		case reflect.Int32:
			val, err := strconv.ParseInt(string(d), 10, 32)
			if err == nil {
				field.Set(reflect.ValueOf(val))
			}
		case reflect.Int64:
			val, err := strconv.ParseInt(string(d), 10, 64)
			if err == nil {
				field.Set(reflect.ValueOf(val))
			}

		case reflect.Float32:
			val, err := strconv.ParseFloat(string(d), 32)
			if err == nil {
				field.Set(reflect.ValueOf(float32(val)))
			}

		case reflect.Float64:
			val, err := strconv.ParseFloat(string(d), 64)
			if err == nil {
				field.Set(reflect.ValueOf(val))
			}

		case reflect.Bool:
			strVal := string(d)
			val := strings.ToLower(strVal) == "true" || strVal == "1"
			field.Set(reflect.ValueOf(val))
			return

		//接口类型
		case reflect.Struct:
			//fmt.Println(reflect.TypeOf(time.Now()), field.Type())
			if reflect.TypeOf(time.Now()) == field.Type() {
				t, err := time.Parse("2006-01-02 15:04:05", string(d))
				if err == nil {
					field.Set(reflect.ValueOf(t.Local()))
				}
			}
			return
		}

	}
}

//遍历所有列，并得到参数及列名
func ItrFieldForSave(meta *TableMeta, val *reflect.Value, includePk bool) (params []interface{}, fieldArr []string) {
	var isSet bool
	for i, k := range meta.FieldMapNames {

		if !includePk && meta.PkIsAuto &&
			meta.FieldMapNames[i] == meta.PkFieldName {
			continue
		}

		field := val.Field(meta.FieldsIndex[i]) // 获取字段所在定义中的位置
		isSet = false

		switch field.Type().Kind() {
		case reflect.String:
			if field.String() != "" {
				isSet = true
				if val.Kind() == reflect.Ptr {
					params = append(params, field.String())
				} else {
					params = append(params, field.String())
				}
			}
			break
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			//if field.Int() != 0 {
			isSet = true
			params = append(params, field.Int())
		//}
		case reflect.Float32, reflect.Float64:
			//if v := field.Float(); v != 0 {
			isSet = true
			params = append(params, field.Float())
		//}

		case reflect.Bool:
			strVal := field.String()
			val := strings.ToLower(strVal) == "true" || strVal == "1"
			field.Set(reflect.ValueOf(val))
			break

		case reflect.Struct:
			v := field.Interface()
			switch v.(type) {
			case time.Time:
				if v.(time.Time).Year() > 1 {
					isSet = true
					params = append(params, v.(time.Time))
				}
			}
		}

		if isSet {
			fieldArr = append(fieldArr, k)
		}
	}
	return params, fieldArr
}

func ItrField(meta *TableMeta, val *reflect.Value, includePk bool) (params []interface{}, fieldArr []string) {
	var isSet bool
	for i, k := range meta.FieldMapNames {

		if !includePk && meta.PkIsAuto &&
			meta.FieldMapNames[i] == meta.PkFieldName {
			continue
		}

		field := val.Field(i)
		isSet = false

		switch field.Type().Kind() {
		case reflect.String:
			if field.String() != "" {
				isSet = true
				if val.Kind() == reflect.Ptr {
					params = append(params, field.String())
				} else {
					params = append(params, field.String())
				}
			}
			break
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if field.Int() != 0 {
				isSet = true
				params = append(params, field.Int())
			}
		case reflect.Float32, reflect.Float64:
			if v := field.Float(); v != 0 {
				isSet = true
				params = append(params, field.Float())
			}

		//		case reflect.Bool:
		//			val := strings.ToLower(strVal) == "true" || strVal == "1"
		//			field.Set(reflect.ValueOf(val))
		//			break

		case reflect.Struct:
			v := field.Interface()
			switch v.(type) {
			case time.Time:
				if v.(time.Time).Year() > 1 {
					isSet = true
					params = append(params, v.(time.Time))
				}
			}
		}

		if isSet {
			fieldArr = append(fieldArr, k)
		}
	}
	return params, fieldArr
}

//************  HELPER  ************//

// save entity and return pk and error
func Save(o Orm, entity interface{}, pk int) (returnPk int, err error) {
	if pk > 0 {
		_, _, err = o.Save(pk, entity)
		return pk, err
	}
	var id64 int64
	_, id64, err = o.Save(nil, entity)
	return int(id64), err
}
