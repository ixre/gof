package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ixre/gof/util"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type (
	TableMapMeta struct {
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

		Dialect() Dialect

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

		//todo:??? 去掉primary参数，并默认Update，如果无返回且无错。则Insert
		Save(primary interface{}, entity interface{}) (rows int64, lastInsertId int64, err error)
	}

	// 表
	Table struct {
		Name    string
		Comment string
		Engine  string
		Charset string
		Columns []*Column
	}

	// 列
	Column struct {
		Name    string
		Pk      bool
		Auto    bool
		NotNull bool
		Type    string
		Comment string
	}

	// find some information of entity
	OrmFinder interface {
	}

	Dialect interface {
		// 数据库方言名称
		Name() string
		// 获取所有的表
		Tables(db *sql.DB, dbName string) ([]*Table, error)
		// 获取表结构
		Table(db *sql.DB, table string) (*Table, error)
	}
)

// 获取表元数据
func GetTableMapMeta(driver string, t reflect.Type) *TableMapMeta {
	ixs, maps := GetFields(driver, t)
	pkName, pkIsAuto := GetPKName(t)
	m := &TableMapMeta{
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
func GetFields(driver string, t reflect.Type) (posArr []int, mapNames []string) {
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
		internalKeysCheck(driver, &fmn)
		mapNames = append(mapNames, fmn)
		posArr = append(posArr, i)
		fmn = ""
	}
	return posArr, mapNames
}

// format internal keywords
func internalKeysCheck(driver string, field *string) {
	if driver == "mysql" {
		checkMysqlInternalKeys(field)
	}
}

func checkMysqlInternalKeys(field *string) {
	switch *field {
	case "key", "where", "type", "describe":
		*field = strings.Join([]string{"`", *field, "`"}, "")
	}
}

func assignValues(meta *TableMapMeta, dst *reflect.Value, rawBytes [][]byte) error {
	for i, fi := range meta.FieldsIndex {
		assignValue(dst.Field(fi), rawBytes[i])
	}
	return nil
}

func assignValue(d reflect.Value, s []byte) (err error) {
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
		v := strings.ToLower(string(s))
		d.SetBool(v == "true" || v == "on" || v == "1")
	case reflect.String:
		d.SetString(string(s))
	case reflect.Slice:
		if d.Type().Elem().Kind() != reflect.Uint8 {
			err = errors.New(fmt.Sprintf("can't covert %s to slice!",
				reflect.TypeOf(s).String()))
		} else {
			d.SetBytes(s)
		}
	}
	return err
}

//遍历所有列，并得到参数及列名
func ItrFieldForSave(meta *TableMapMeta, val *reflect.Value, includePk bool) (
	params []interface{}, fieldArr []string) {
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
			isSet = true
			if val.Kind() == reflect.Ptr {
				params = append(params, field.String())
			} else {
				params = append(params, field.String())
			}
		case reflect.Int, reflect.Int8,
			reflect.Int16, reflect.Int32, reflect.Int64:
			isSet = true
			params = append(params, field.Int())

		case reflect.Float32, reflect.Float64:
			isSet = true
			params = append(params, field.Float())

		case reflect.Bool:
			strVal := field.String()
			val := strings.ToLower(strVal) == "true" || strVal == "1"
			field.Set(reflect.ValueOf(val))

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

// save entity and return pk and error
func Save(o Orm, entity interface{}, pk int) (int, error) {
	if pk > 0 {
		_, _, err := o.Save(pk, entity)
		return pk, err
		//r, _, err := o.Save(pk, entity)
		//if r == 0 && err == nil {
		//    _, _, err = o.Save(nil, entity)
		//}
		//return pk, err
	}
	_, int64, err := o.Save(nil, entity)
	return int(int64), err
}

// parse save result int to int32
func I32(v int, err error) (int32, error) {
	return util.I32Err(v, err)
}

func I64(v int, err error) (int64, error) {
	return util.I64Err(v, err)
}
