/**
 * Copyright 2015 @ at3.net.
 * name : mssql_dialect
 * author : jarryliu
 * date : 2016-11-11 12:29
 * description :
 * history :
 */
package dialect

import (
	"database/sql"
	"fmt"

	"github.com/ixre/gof/db/db"
	"github.com/ixre/gof/types/typeconv"
)

var _ Dialect = new(MsSqlDialect)

type MsSqlDialect struct {
}

func (m *MsSqlDialect) GetField(f string) string {
	return fmt.Sprintf("[%s]", f)
}

func (m *MsSqlDialect) Name() string {
	return "MSSQLDialect"
}

// 获取所有的表
func (m *MsSqlDialect) Tables(d *sql.DB, dbName string, schema string) ([]*db.Table, error) {
	var list []string
	stmt, err := d.Prepare(`
	 SELECT top 1 ob.name FROM sys.objects AS ob
      LEFT OUTER JOIN sys.extended_properties AS ep
        ON ep.major_id = ob.object_id
           AND ep.class = 1
           AND ep.minor_id = 0
    WHERE ObjectProperty(ob.object_id, 'IsUserTable') = 1 `)
	if err == nil {
		tb := ""
		rows, _ := stmt.Query()
		for rows.Next() {
			if err = rows.Scan(&tb); err == nil {
				list = append(list, tb)
			}
		}
		stmt.Close()
		rows.Close()
		tList := make([]*db.Table, len(list))
		for i, v := range list {
			if tList[i], err = m.Table(d, v); err != nil {
				return nil, err
			}
		}
		return tList, nil
	}
	return nil, err
}

// 获取表结构
func (m *MsSqlDialect) Table(d *sql.DB, table string) (*db.Table, error) {
	stmt, err := d.Prepare("exec sp_columns " + table)
	if err == nil {
		table := &db.Table{
			Name:    table,
			Comment: "",
			Engine:  "",
			Charset: "",
			Columns: []*db.Column{},
		}
		if table.Comment != "" {
			table.Comment = ""
		}

		rs := make([]interface{}, 19)
		var rawBytes = make([][]byte, 19)
		for i := range rs {
			rs[i] = &rawBytes[i]
		}
		rows, _ := stmt.Query()
		for rows.Next() {
			err = rows.Scan(rs...)
			if err == nil {
				table.Columns = append(table.Columns, m.parseColumn(rs))
			}
		}
		stmt.Close()
		m.updatePkColumn(d,table)
		m.updateAutoKeys(d,table)
		return table,nil
	}
	return nil, err
}

// 更新主键
func (m *MsSqlDialect) updatePkColumn(d *sql.DB,table *db.Table){
	stmt, _ := d.Prepare(fmt.Sprintf(`exec sp_pkeys %s`,table.Name))
	rows, _ := stmt.Query()
	rs := make([]interface{}, 6)
		var rawBytes = make([][]byte, 6)
		for i := range rs {
			rs[i] = &rawBytes[i]
	}
	if rows.Next(){
		_ = rows.Scan(rs...)
		pk := getString(rs,3)
		for _,v := range table.Columns{
			if v.Name == pk{
				v.IsPk = true
				break
			}
		}
	}
	stmt.Close()
}

// 更新自增键
func (m *MsSqlDialect) updateAutoKeys(d *sql.DB,table *db.Table){
	stmt, err := d.Prepare(fmt.Sprintf(`
	SELECT name FROM syscolumns   
	WHERE id=object_id('%s') 
	AND COLUMNPROPERTY(id,name,'IsIdentity')=1`,table.Name))
	rows, _ := stmt.Query()
	keys := make(map[string]int,0)
	for rows.Next() {
		s := ""
		err = rows.Scan(&s)
		if err == nil {
			keys[s] = 1
		}
	}
	stmt.Close()
	for _,v := range table.Columns{
		if _,ok := keys[v.Name];ok {
			v.IsAuto = true
		}
	}
}

func (m *MsSqlDialect) getTypeId(dbType string) int {
	switch dbType {
	case "smallint":
		return db.TypeInt16
	case "tinyint", "int":
		return db.TypeInt32
	case "bit":
		return db.TypeBoolean
	case "bigint":
		return db.TypeInt64
	case "float":
		return db.TypeFloat32
	case "decimal":
		return db.TypeDecimal
	case "double":
		return db.TypeFloat64
	case "date":
		return db.TypeDateTime
	case "text", "varchar", "char":
		return db.TypeString
	}
	println("[ ORM][ MSSQL][ Warning]:Dialect not support type :", dbType)
	return db.TypeUnknown
}

func getBytes(rs []interface{}, i int) []byte {
	return *(rs[i].(*[]byte))
}

func getString(rs []interface{}, i int) string {
	s, _ := typeconv.String(getBytes(rs, i))
	return s
}

func (m *MsSqlDialect) parseColumn(rs []interface{}) *db.Column {
	/*
		--获取某表中的自动增长列的列名
		select   name   from   syscolumns
		  where   id=object_id('corpinfo')   and
		                COLUMNPROPERTY(id,name,'IsIdentity')=1
	*/

	dbType := getString(rs, 5)
	return &db.Column{
		Name:    getString(rs, 3),
		IsPk:    false,
		NotNull: !typeconv.MustBool(getString(rs, 10)),
		IsAuto:  false,
		DbType:  dbType,
		Length:  typeconv.MustInt(getString(rs, 7)),
		Type:    m.getTypeId(dbType),
		Comment: getString(rs, 11),
	}
}
