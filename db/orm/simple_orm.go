package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jsix/gof/storage"
	"log"
	"reflect"
	"strings"
)

var _ Orm = new(simpleOrm)

//it's a IOrm Implements for mysql
type simpleOrm struct {
	tableMap map[string]*TableMapMeta
	*sql.DB
	driverName string
	useTrace   bool
	dialect    Dialect
}

func NewOrm(driver string, db *sql.DB) Orm {
	var dialect Dialect
	switch driver {
	case "mysql":
		dialect = &MySqlDialect{}
	case "mssql":
		dialect = &MsSqlDialect{}
	case "sqlite":
		dialect = &SqLiteDialect{}
	}
	return &simpleOrm{
		DB:         db,
		driverName: driver,
		dialect:    dialect,
		tableMap:   make(map[string]*TableMapMeta),
	}
}

func NewCachedOrm(driver string, db *sql.DB, s storage.Interface) Orm {
	return NewCacheProxyOrm(NewOrm(driver, db), s)
}

func (o *simpleOrm) Version() string {
	return "1.0.2"
}

func (s *simpleOrm) Dialect() Dialect {
	return s.dialect
}

func (o *simpleOrm) err(err error) error {
	if o.useTrace && err != nil && err != sql.ErrNoRows {
		log.Println("[ ORM][ ERROR]:", err.Error())
	}
	return err
}

func (o *simpleOrm) debug(format string, args ...interface{}) {
	if o.useTrace {
		log.Printf(format+"\n", args...)
	}
}

func (o *simpleOrm) getTableMapMeta(t reflect.Type) *TableMapMeta {
	m, exists := o.tableMap[t.String()]
	if exists {
		return m
	}
	m = GetTableMapMeta(o.driverName, t)
	o.tableMap[t.String()] = m

	if o.useTrace {
		log.Println("[ ORM][ META]:", m)
	}

	return m
}

func (o *simpleOrm) getTableName(t reflect.Type) string {
	//todo: 用int做键
	v, exists := o.tableMap[t.String()]
	if exists {
		return v.TableName
	}
	return t.Name()
}

//if not defined primary key.the first key will as primary key
func (o *simpleOrm) getPKName(t reflect.Type) (pkName string, pkIsAuto bool) {
	v, exists := o.tableMap[t.String()]
	if exists {
		return v.PkFieldName, v.PkIsAuto
	}
	return GetPKName(t)
}

func (o *simpleOrm) unionField(meta *TableMapMeta, v string) string {
	if len(meta.TableName) != 0 {
		return meta.TableName + "." + v
	}
	return v
}

func (o *simpleOrm) SetTrace(b bool) {
	o.useTrace = b
}

//create a fixed table map
func (o *simpleOrm) Mapping(v interface{}, table string) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	meta := o.getTableMapMeta(t)
	meta.TableName = table
	o.tableMap[t.String()] = meta
}

func (o *simpleOrm) Get(primaryVal interface{}, entity interface{}) error {
	var sql string
	var fieldLen int
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	val := reflect.ValueOf(entity)
	if val.Kind() != reflect.Ptr {
		return o.err(errors.New("Unaddressable of entity ,it must be a ptr"))
	}
	val = val.Elem()
	/* build sql */
	meta := o.getTableMapMeta(t)
	fieldLen = len(meta.FieldsIndex)
	fieldArr := make([]string, fieldLen)
	var scanVal []interface{} = make([]interface{}, fieldLen)
	var rawBytes [][]byte = make([][]byte, fieldLen)
	for i, v := range meta.FieldMapNames {
		fieldArr[i] = v
		scanVal[i] = &rawBytes[i]
	}
	sql = fmt.Sprintf("SELECT %s FROM %s WHERE %s=?",
		strings.Join(fieldArr, ","),
		meta.TableName,
		meta.PkFieldName,
	)
	if o.useTrace {
		log.Println(fmt.Sprintf("[ ORM][ SQL]:%s , [ Params]:%+v", sql, primaryVal))
	}
	/* query */
	stmt, err := o.DB.Prepare(sql)
	if err != nil {
		return o.err(errors.New(err.Error() + "\n[ SQL]:" + sql))
	}
	defer stmt.Close()
	row := stmt.QueryRow(primaryVal)
	err = row.Scan(scanVal...)
	if err != nil {
		return o.err(err)
	}
	return BindFields(meta, &val, rawBytes)
}

func (o *simpleOrm) GetBy(entity interface{}, where string,
	args ...interface{}) error {

	var sql string
	var fieldLen int
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	val := reflect.ValueOf(entity)
	if val.Kind() != reflect.Ptr {
		return o.err(errors.New("unaddressable of entity ,it must be a ptr"))
	}

	if strings.Trim(where, "") == "" {
		return o.err(errors.New("param where can't be empty "))
	}

	val = val.Elem()

	if !val.IsValid() {
		return o.err(errors.New("not validate or not initialize."))
	}

	/* build sql */
	meta := o.getTableMapMeta(t)
	fieldLen = len(meta.FieldsIndex)
	fieldArr := make([]string, fieldLen)

	var scanVal []interface{} = make([]interface{}, fieldLen)
	var rawBytes [][]byte = make([][]byte, fieldLen)

	for i, v := range meta.FieldMapNames {
		fieldArr[i] = v
		scanVal[i] = &rawBytes[i]
	}

	sql = fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(fieldArr, ","),
		meta.TableName,
		where,
	)

	if o.useTrace {
		log.Println(fmt.Sprintf("[ ORM][ SQL]:%s , [ Params]:%s - %+v", sql, where, args))
	}

	/* query */
	stmt, err := o.DB.Prepare(sql)
	if err != nil {
		if o.useTrace {
			log.Println("[ ORM][ ERROR]:", err.Error(), " [ SQL]:", sql)
		}
		return o.err(errors.New(err.Error() + "\n[ SQL]:" + sql))
	}
	defer stmt.Close()

	row := stmt.QueryRow(args...)
	err = row.Scan(scanVal...)

	if err != nil {
		return o.err(err)
	}

	return BindFields(meta, &val, rawBytes)
}

func (o *simpleOrm) GetByQuery(entity interface{}, sql string,
	args ...interface{}) error {
	var fieldLen int
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	val := reflect.ValueOf(entity)
	if val.Kind() != reflect.Ptr {
		return o.err(errors.New("Unaddressable of entity ,it must be a ptr"))
	}

	val = val.Elem()

	/* build sql */
	meta := o.getTableMapMeta(t)
	fieldLen = len(meta.FieldsIndex)
	fieldArr := make([]string, fieldLen)
	var scanVal []interface{} = make([]interface{}, fieldLen)
	var rawBytes [][]byte = make([][]byte, fieldLen)

	for i, v := range meta.FieldMapNames {
		fieldArr[i] = o.unionField(meta, v)
		scanVal[i] = &rawBytes[i]
	}

	if strings.Index(sql, "*") != -1 {
		sql = strings.Replace(sql, "*", strings.Join(fieldArr, ","), 1)
	}

	if o.useTrace {
		log.Println(fmt.Sprintf("[ ORM][ SQL]:%s", sql))
	}

	/* query */
	stmt, err := o.DB.Prepare(sql)

	if err != nil {
		return o.err(errors.New(err.Error() + "\n[ SQL]:" + sql))
	}
	defer stmt.Close()

	row := stmt.QueryRow(args...)
	err = row.Scan(scanVal...)

	if err != nil {
		return o.err(err)
	}

	return BindFields(meta, &val, rawBytes)
}

//Select more than 1 entity list
//@to : referenced queried entity list
//@entity : query condition
//@where : other condition
func (o *simpleOrm) Select(to interface{}, where string, args ...interface{}) error {
	return o.selectBy(to, where, false, args...)
}

func (o *simpleOrm) SelectByQuery(to interface{}, sql string, args ...interface{}) error {
	return o.selectBy(to, sql, true, args...)
}

// query rows
func (o *simpleOrm) selectBy(dst interface{}, sql string, fullSql bool, args ...interface{}) error {
	var fieldLen int
	var eleIsPtr bool // 元素是否为指针

	toVal := reflect.Indirect(reflect.ValueOf(dst))
	toTyp := reflect.TypeOf(dst).Elem()

	if toTyp.Kind() == reflect.Ptr {
		toTyp = toTyp.Elem()
	}

	if toTyp.Kind() != reflect.Slice {
		return o.err(errors.New("dst must be slice"))
	}

	baseType := toTyp.Elem()
	if baseType.Kind() == reflect.Ptr {
		eleIsPtr = true
		baseType = baseType.Elem()
	}

	/* build sql */
	meta := o.getTableMapMeta(baseType)
	fieldLen = len(meta.FieldMapNames)
	fieldArr := make([]string, fieldLen)
	var scanVal []interface{} = make([]interface{}, fieldLen)
	var rawBytes [][]byte = make([][]byte, fieldLen)

	for i, v := range meta.FieldMapNames {
		fieldArr[i] = o.unionField(meta, v)
		scanVal[i] = &rawBytes[i]
	}

	if fullSql {
		if strings.Index(sql, "*") != -1 {
			sql = strings.Replace(sql, "*", strings.Join(fieldArr, ","), 1)
		}
	} else {
		where := sql
		if len(where) == 0 {
			sql = fmt.Sprintf("SELECT %s FROM %s",
				strings.Join(fieldArr, ","),
				meta.TableName)
		} else {
			// 此时,sql为查询条件
			sql = fmt.Sprintf("SELECT %s FROM %s WHERE %s",
				strings.Join(fieldArr, ","),
				meta.TableName,
				where)
		}
	}

	if o.useTrace {
		log.Println(fmt.Sprintf("[ ORM][ SQL]:%s [ Params] - %+v", sql, args))
	}

	/* query */
	stmt, err := o.DB.Prepare(sql)
	if err != nil {
		return o.err(errors.New(fmt.Sprintf("%s - [ SQL]: %s- [Args]:%+v", err.Error(), sql, args)))
	}

	defer stmt.Close()
	rows, err := stmt.Query(args...)

	if err != nil {
		return o.err(errors.New(err.Error() + "\n[ SQL]:" + sql))
	}

	defer rows.Close()

	/* 用反射来对输出结果复制 */

	toArr := toVal

	for rows.Next() {
		e := reflect.New(baseType)
		v := e.Elem()
		if err = rows.Scan(scanVal...); err != nil {
			break
		}
		//for i, fi := range meta.FieldsIndex {
		//	SetField(v.Field(fi), rawBytes[i])
		//}

		BindFields(meta, &v, rawBytes)
		if eleIsPtr {
			toArr = reflect.Append(toArr, e)
		} else {
			toArr = reflect.Append(toArr, v)
		}
	}
	toVal.Set(toArr)
	return o.err(err)
}

func (o *simpleOrm) Delete(entity interface{}, where string,
	args ...interface{}) (effect int64, err error) {
	var sql string

	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	/* build sql */
	meta := o.getTableMapMeta(t)

	if where == "" {
		return 0, errors.New("Unknown condition")
	}

	sql = fmt.Sprintf("DELETE FROM %s WHERE %s",
		meta.TableName,
		where,
	)

	if o.useTrace {
		log.Println(fmt.Sprintf("[ ORM][ SQL]:%s , [ Params]:%#v", sql, args))
	}

	/* query */
	stmt, err := o.DB.Prepare(sql)
	if err != nil {
		return 0, o.err(errors.New(fmt.Sprintf("[ ORM][ ERROR]:%s [ SQL]:%s", err.Error(), sql)))

	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	var rowNum int64 = 0
	if err == nil {
		rowNum, err = result.RowsAffected()
	}
	if err != nil {
		return rowNum, o.err(errors.New(err.Error() + "\n[ SQL]" + sql))
	}
	return rowNum, nil
}

func (o *simpleOrm) DeleteByPk(entity interface{}, primary interface{}) (err error) {
	var sql string
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	/* build sql */
	meta := o.getTableMapMeta(t)

	sql = fmt.Sprintf("DELETE FROM %s WHERE %s=?",
		meta.TableName,
		meta.PkFieldName,
	)

	if o.useTrace {
		log.Println(fmt.Sprintf("[ ORM][ SQL]:%s , [ Params]:%s", sql, primary))
	}

	/* query */
	stmt, err := o.DB.Prepare(sql)
	if err != nil {
		return o.err(errors.New(fmt.Sprintf("[ ORM][ ERROR]:%s \n [ SQL]:%s", err.Error(), sql)))

	}
	defer stmt.Close()

	_, err = stmt.Exec(primary)
	if err != nil {
		return o.err(errors.New(err.Error() + "\n[ SQL]" + sql))
	}
	return nil
}

func (o *simpleOrm) Save(primaryKey interface{}, entity interface{}) (rows int64, lastInsertId int64, err error) {
	var sql string
	//var condition string
	//var fieldLen int

	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	val := reflect.Indirect(reflect.ValueOf(entity))

	// build sql
	meta := o.getTableMapMeta(t)
	//fieldLen = len(meta.FieldNames)
	params, fieldArr := ItrFieldForSave(meta, &val, false)

	//insert
	if primaryKey == nil {
		var pArr = make([]string, len(fieldArr))
		for i := range pArr {
			pArr[i] = "?"
		}

		sql = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", meta.TableName,
			strings.Join(fieldArr, ","),
			strings.Join(pArr, ","),
		)

		if o.useTrace {
			log.Println(fmt.Sprintf("[ ORM][ SQL]:%s , [ Params]:%+v", sql, params))
		}

		/* query */
		stmt, err := o.DB.Prepare(sql)
		if err != nil {
			return 0, 0, o.err(errors.New("[ ORM][ ERROR]:" + err.Error() + "\n[ SQL]" + sql))
		}
		defer stmt.Close()

		result, err := stmt.Exec(params...)
		var rowNum int64 = 0
		var lastInsertId int64 = 0
		if err == nil {
			rowNum, err = result.RowsAffected()
			if err == nil {
				lastInsertId, err = result.LastInsertId()
			}
			return rowNum, lastInsertId, err
		}
		return rowNum, lastInsertId, o.err(errors.New(err.Error() + "\n[ SQL]" + sql))
	} else {
		//update model

		var setCond string

		for i, k := range fieldArr {
			if i == 0 {
				setCond = fmt.Sprintf("%s = ?", k)
			} else {
				setCond = fmt.Sprintf("%s,%s = ?", setCond, k)
			}
		}

		sql = fmt.Sprintf("UPDATE %s SET %s WHERE %s=?", meta.TableName,
			setCond,
			meta.PkFieldName,
		)

		/* query */
		stmt, err := o.DB.Prepare(sql)
		if err != nil {
			return 0, 0, o.err(errors.New("[ ORM][ ERROR]:" + err.Error() + " [ SQL]" + sql))
		}
		defer stmt.Close()

		params = append(params, primaryKey)

		if o.useTrace {
			log.Println(fmt.Sprintf("[ ORM][ SQL]:%s , [ Params]:%+v", sql, params))
		}

		result, err := stmt.Exec(params...)
		var rowNum int64 = 0
		if err == nil {
			rowNum, err = result.RowsAffected()
			return rowNum, 0, err
		}
		return rowNum, 0, o.err(errors.New(err.Error() + "\n[ SQL]" + sql))
	}
}
