package orm

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ixre/gof/storage"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var _ Orm = new(simpleOrm)

//it's a IOrm Implements for mysql
type simpleOrm struct {
	tableMap map[string]*TableMapMeta
	*sql.DB
	driverName string
	useTrace   bool
	dialect    Dialect
	tmLock     *sync.RWMutex
}

func NewOrm(driver string, db *sql.DB) Orm {
	var dialect Dialect
	switch driver {
	case "mysql":
		dialect = &MySqlDialect{}
	case "postgres", "postgresql":
		dialect = &PostgresqlDialect{}
	case "mssql":
		dialect = &MsSqlDialect{}
	case "sqlite":
		dialect = &SqLiteDialect{}
	}
	return &simpleOrm{
		DB:         db,
		driverName: strings.ToLower(driver),
		dialect:    dialect,
		tableMap:   make(map[string]*TableMapMeta),
		tmLock:     &sync.RWMutex{},
	}
}

func NewCachedOrm(driver string, db *sql.DB, s storage.Interface) Orm {
	return NewCacheProxyOrm(NewOrm(driver, db), s)
}

func (o *simpleOrm) Version() string {
	return "1.0.2"
}

func (o *simpleOrm) Dialect() Dialect {
	return o.dialect
}

func (o *simpleOrm) err(err error, s string, args ...interface{}) error {
	if err != nil && err != sql.ErrNoRows {
		if len(s) == 0 {
			log.Println("[ ORM][ ERROR]:", err.Error())
		} else {
			if len(args) > 0 {
				log.Println(fmt.Sprintf("[ ORM][ ERROR]:%s [ SQL]:%s [Args]:%+v", err.Error(), s, args))
			} else {
				log.Println(fmt.Sprintf("[ ORM][ ERROR]:%s [ SQL]:%s ", err.Error(), s))
			}
		}
	}
	return err
}

func (o *simpleOrm) debug(format string, args ...interface{}) {
	if o.useTrace {
		log.Printf(format+"\n", args...)
	}
}

func (o *simpleOrm) getTableMapMeta(t reflect.Type) *TableMapMeta {
	o.tmLock.RLock()
	m, exists := o.tableMap[t.String()]
	o.tmLock.RUnlock()
	if exists {
		return m
	}
	o.tmLock.Lock()
	m = GetTableMapMeta(o.driverName, t)
	o.tableMap[t.String()] = m
	o.tmLock.Unlock()
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

func (o *simpleOrm) unionField(meta *TableMapMeta, field string) string {
	if len(meta.TableName) != 0 {
		return meta.TableName + "." + field
	}
	return field
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
	entityName := t.String()
	m := o.getTableMapMeta(t)
	m.TableName = table
	o.tableMap[entityName] = m
	if o.useTrace {
		logTxt := fmt.Sprintf("TableName=%s,FieldCount=%d,PkFieldName=%s,PkIsAuto=%v",
			m.TableName, len(m.FieldMapNames), m.PkFieldName, m.PkIsAuto)
		log.Println("[ ORM][ Mapping]:", entityName, "->", logTxt)
	}
}

func (o *simpleOrm) Get(primaryVal interface{}, entity interface{}) error {
	var fieldLen int
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	val := reflect.ValueOf(entity)
	if val.Kind() != reflect.Ptr {
		return o.err(errors.New("Unaddressable of entity ,it must be a ptr"), "")
	}
	val = val.Elem()
	/* build sqlQuery */
	meta := o.getTableMapMeta(t)
	fieldLen = len(meta.FieldsIndex)
	fieldArr := make([]string, fieldLen)
	var scanVal = make([]interface{}, fieldLen)
	var rawBytes = make([][]byte, fieldLen)
	for i, v := range meta.FieldMapNames {
		fieldArr[i] = v
		scanVal[i] = &rawBytes[i]
	}
	sqlQuery := o.fmtSelectSingleQuery(fieldArr, meta.TableName, meta.PkFieldName+" = "+o.getParamHolder(0))
	if o.useTrace {
		log.Println(fmt.Sprintf("[ ORM][ SQL]:%s , [ Params]:%+v", sqlQuery, primaryVal))
	}
	stmt, err := o.DB.Prepare(sqlQuery)
	if err == nil {
		defer stmt.Close()
		row := stmt.QueryRow(primaryVal)
		err = row.Scan(scanVal...)
	}
	if err != nil {
		return o.err(err, sqlQuery)
	}
	return assignValues(meta, &val, rawBytes)
}

func (o *simpleOrm) GetBy(entity interface{}, where string,
	args ...interface{}) error {

	var fieldLen int
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	val := reflect.ValueOf(entity)
	if val.Kind() != reflect.Ptr {
		return o.err(errors.New("unaddressable of entity ,it must be a ptr"), "")
	}

	if strings.Trim(where, "") == "" {
		return o.err(errors.New("param where can't be empty "), "")
	}

	val = val.Elem()

	if !val.IsValid() {
		return o.err(errors.New("not validate or not initialize."), "")
	}

	/* build sqlQuery */
	meta := o.getTableMapMeta(t)
	fieldLen = len(meta.FieldsIndex)
	fieldArr := make([]string, fieldLen)

	var scanVal = make([]interface{}, fieldLen)
	var rawBytes = make([][]byte, fieldLen)

	for i, v := range meta.FieldMapNames {
		fieldArr[i] = v
		scanVal[i] = &rawBytes[i]
	}
	sqlQuery := o.fmtSelectSingleQuery(fieldArr, meta.TableName, where)
	if o.useTrace {
		log.Println(fmt.Sprintf("[ ORM][ SQL]:%s , [ Params]:%s - %+v", sqlQuery, where, args))
	}
	/* query */
	stmt, err := o.DB.Prepare(sqlQuery)
	if err == nil {
		defer stmt.Close()
		row := stmt.QueryRow(args...)
		err = row.Scan(scanVal...)
	}
	if err != nil {
		return o.err(err, sqlQuery)
	}
	return assignValues(meta, &val, rawBytes)
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
		return o.err(errors.New("Unaddressable of entity ,it must be a ptr"), "")
	}

	val = val.Elem()
	/* build sql */
	meta := o.getTableMapMeta(t)
	fieldLen = len(meta.FieldsIndex)
	fieldArr := make([]string, fieldLen)
	var scanVal = make([]interface{}, fieldLen)
	var rawBytes = make([][]byte, fieldLen)
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
	if err == nil {
		defer stmt.Close()
		row := stmt.QueryRow(args...)
		err = row.Scan(scanVal...)
	}
	if err != nil {
		return o.err(err, sql)
	}
	return assignValues(meta, &val, rawBytes)
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
		return o.err(errors.New("dst must be slice"), "")
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
	var scanVal = make([]interface{}, fieldLen)
	var rawBytes = make([][]byte, fieldLen)
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
		return o.err(err, sql, args)
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return o.err(err, sql)
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
		//println(fmt.Sprintf("%#v",	string(rawBytes[0])))
		//println(fmt.Sprintf("%#v",	string(rawBytes[1])))
		//for i, fi := range meta.FieldsIndex {
		//	SetField(v.Field(fi), rawBytes[i])
		//}
		assignValues(meta, &v, rawBytes)
		if eleIsPtr {
			toArr = reflect.Append(toArr, e)
		} else {
			toArr = reflect.Append(toArr, v)
		}
	}
	toVal.Set(toArr)
	return o.err(err, sql)
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
		log.Println(fmt.Sprintf("[ ORM][ SQL]:%s , [ Params]:%+v", sql, args))
	}
	/* query */
	stmt, err := o.DB.Prepare(sql)
	if err != nil {
		return 0, o.err(err, sql)
	}
	defer stmt.Close()
	result, err := stmt.Exec(args...)
	var rowNum int64 = 0
	if err == nil {
		rowNum, err = result.RowsAffected()
	}
	if err != nil {
		return rowNum, o.err(err, sql)
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

	sql = fmt.Sprintf("DELETE FROM %s WHERE %s=%s",
		meta.TableName,
		meta.PkFieldName,
		o.getParamHolder(0),
	)

	if o.useTrace {
		log.Println(fmt.Sprintf("[ ORM][ SQL]:%s , [ Params]:%+v", sql, primary))
	}
	/* query */
	stmt, err := o.DB.Prepare(sql)
	if err != nil {
		return o.err(err, sql)
	}
	defer stmt.Close()
	_, err = stmt.Exec(primary)
	if err != nil {
		return o.err(err, sql)
	}
	return nil
}

func (o *simpleOrm) Save(primary interface{}, entity interface{}) (rows int64, lastInsertId int64, err error) {
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
	// pk type is int?
	isIntPk := o.isIntPk(meta.PkFieldTypeId)
	//fieldLen = len(meta.FieldNames)
	params, fieldArr := ItrFieldForSave(meta, &val, false)

	//insert
	if primary == nil {
		var pArr = make([]string, len(fieldArr))
		for i := range pArr {
			pArr[i] = o.getParamHolder(i)
		}
		sql = o.getInsertSQL(meta, fieldArr, pArr)
		if o.useTrace {
			log.Println(fmt.Sprintf("[ ORM][ SQL]:%s , [ Params]:%+v", sql, params))
		}
		/* query */
		stmt, err := o.DB.Prepare(sql)
		if err != nil {
			return 0, 0, o.err(err, sql)
		}
		defer stmt.Close()
		return o.stmtUpdateExec(isIntPk, stmt, sql, params...)
	} else {
		//update model
		var setCond string
		for i, k := range fieldArr {
			if i == 0 {
				setCond = fmt.Sprintf("%s = %s", k, o.getParamHolder(i))
			} else {
				setCond = fmt.Sprintf("%s,%s = %s", setCond, k, o.getParamHolder(i))
			}
		}

		sql = o.getUpdateSQL(meta, setCond, fieldArr)
		stmt, err := o.DB.Prepare(sql)
		if err != nil {
			return 0, 0, o.err(err, sql)
		}
		defer stmt.Close()
		params = append(params, primary)
		if o.useTrace {
			log.Println(fmt.Sprintf("[ ORM][ SQL]:%s , [ Params]:%+v", sql, params))
		}
		return o.stmtUpdateExec(isIntPk, stmt, sql, params...)
	}
}

func (o *simpleOrm) stmtUpdateExec(isIntPk bool, stmt *sql.Stmt, sql_ string, params ...interface{}) (int64, int64, error) {
	// Postgresql 新增或更新时候,使用returning语句,应当做Result查询
	if (o.driverName == "postgres" || o.driverName == "postgresql") && (strings.Contains(sql_, "returning") || strings.Contains(sql_, "RETURNING")) {
		var lastInsertId int64
		row := stmt.QueryRow(params...)
		if isIntPk {
			if err := row.Scan(&lastInsertId); err != nil {
				return 0, lastInsertId, o.err(err, sql_)
			}
		}
		return 0, lastInsertId, nil
	}
	result, err := stmt.Exec(params...)
	var rowNum int64 = 0
	var lastInsertId int64 = 0
	if err == nil {
		rowNum, err = result.RowsAffected()
		if err == nil && isIntPk {
			lastInsertId, err = result.LastInsertId()
		}
		return rowNum, lastInsertId, err
	}
	return rowNum, lastInsertId, o.err(err, sql_)
}

func (o *simpleOrm) fmtSelectSingleQuery(fields []string, table string, where string) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("SELECT ")
	if o.driverName == "mssql" {
		buf.WriteString("TOP 1 ")
	}
	buf.WriteString(strings.Join(fields, ","))
	buf.WriteString(" FROM ")
	buf.WriteString(table)
	if len(where) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(where)
		if o.driverName == "mysql" && strings.Index(
			strings.ToLower(where), " limit ") == -1 {
			buf.WriteString(" LIMIT 1")
		}
	} else {
		//if o.driverName == "postgresql" || o.driverName =="postgres"{
		//	buf.WriteString("LIMIT 1")
		//}
		buf.WriteString(" LIMIT 1")
	}
	return buf.String()
}

func (o *simpleOrm) getParamHolder(n int) string {
	switch o.driverName {
	case "mysql":
		return "?"
	case "postgres", "postgresql":
		return "$" + strconv.Itoa(n+1)
	}
	return "?"
}

// 获取插入数据SQL
func (o *simpleOrm) getInsertSQL(meta *TableMapMeta, fieldArr []string, pArr []string) string {
	buf := bytes.NewBufferString("INSERT INTO ")
	buf.WriteString(meta.TableName)
	buf.WriteString(" (")
	buf.WriteString(strings.Join(fieldArr, ","))
	buf.WriteString(" ) VALUES (")
	buf.WriteString(strings.Join(pArr, ","))
	buf.WriteString(" )")
	// Postgresql需要在INSERT后执行 RETURNING id 才能返回LastInsertId
	if o.driverName == "postgres" || o.driverName == "postgresql" {
		buf.WriteString(" RETURNING ")
		buf.WriteString(meta.PkFieldName)
		buf.WriteString(";")
	}
	return buf.String()
}

// 获取UPDATE SQL
func (o *simpleOrm) getUpdateSQL(meta *TableMapMeta, setCond string, fieldArr []string) string {
	buf := bytes.NewBufferString("UPDATE ")
	buf.WriteString(meta.TableName)
	buf.WriteString(" SET ")
	buf.WriteString(setCond)
	buf.WriteString(" WHERE ")
	buf.WriteString(meta.PkFieldName)
	buf.WriteString(" = ")
	buf.WriteString(o.getParamHolder(len(fieldArr)))
	// Postgresql需要在INSERT后执行 RETURNING id 才能返回LastInsertId
	if o.driverName == "postgres" || o.driverName == "postgresql" {
		buf.WriteString(" RETURNING ")
		buf.WriteString(meta.PkFieldName)
		buf.WriteString(";")
	}
	return buf.String()
}

func (o *simpleOrm) isIntPk(typeId int) bool {
	switch typeId {
	case TypeInt16, TypeInt32, TypeInt64:
		return true
	}
	return false
}
