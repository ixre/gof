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
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/ixre/gof/db/db"
	"github.com/ixre/gof/typeconv"
)

var _ Dialect = new(MsSqlDialect)

type MsSqlDialect struct {
}

func (m *MsSqlDialect) GetField(f string) string {
	if strings.Contains(f, ".") {
		return f
	}
	return fmt.Sprintf("[%s]", f)
}

func (m *MsSqlDialect) Name() string {
	return "MSSQLDialect"
}

func (m *MsSqlDialect) fetchTableNames(d *sql.DB, dbName string, keyword string) (map[string]string, error) {
	buf := bytes.NewBufferString(` SELECT ob.name,ISNULL(ep.value,'') as comment
	  	FROM sys.objects AS ob
	  	LEFT OUTER JOIN sys.extended_properties AS ep
	  	ON ep.major_id = ob.object_id  AND ep.class = 1  AND ep.minor_id = 0
  		WHERE ObjectProperty(ob.object_id, 'IsUserTable') = 1`)
	if keyword != "" {
		buf.WriteString(` AND ob.name LIKE '%`)
		buf.WriteString(keyword)
		buf.WriteString(`%'`)
	}
	var tables = make(map[string]string, 0)
	stmt, err := d.Prepare(buf.String())
	if err == nil {
		tb := ""
		comment := ""
		rows, _ := stmt.Query()
		for rows.Next() {
			if err = rows.Scan(&tb, &comment); err == nil {
				tables[tb] = comment
			}
		}
		stmt.Close()
		rows.Close()
	}
	return tables, err
}

// 获取所有的表
func (m *MsSqlDialect) Tables(d *sql.DB, dbName string, keyword string, match func(int, string) bool) (int, []*db.Table, error) {
	tables, err := m.fetchTableNames(d, dbName, keyword)
	l := len(tables)
	tableList := make([]string, 0, len(tables))
	for k := range tables {
		tableList = append(tableList, k)
	}
	if err == nil {
		i := -1
		tList := make([]*db.Table, 0)
		for _, k := range tableList {
			i++
			if match != nil && !match(i, k) {
				// 筛选掉不匹配的表
				continue
			}
			if tb, err := m.Table(d, k); err == nil {
				tb.Comment = tables[k]
				tList = append(tList, tb)
			} else {
				log.Println("[ mssql][ dialect]: get table structure failed. " + err.Error())
			}
		}
		return l, tList, nil
	}
	return l, nil, err
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
		rows.Close()
		stmt.Close()
		m.updatePkColumn(d, table)
		m.updateColumnOthers(d, table)
		return table, nil
	}
	return nil, err
}

// 更新主键
func (m *MsSqlDialect) updatePkColumn(d *sql.DB, table *db.Table) {
	stmt, _ := d.Prepare(fmt.Sprintf(`exec sp_pkeys %s`, table.Name))
	rows, _ := stmt.Query()
	rs := make([]interface{}, 6)
	var rawBytes = make([][]byte, 6)
	for i := range rs {
		rs[i] = &rawBytes[i]
	}
	if rows.Next() {
		_ = rows.Scan(rs...)
		pk := getString(rs, 3)
		for _, v := range table.Columns {
			if v.Name == pk {
				v.IsPk = true
				break
			}
		}
	}
	rows.Close()
	stmt.Close()
}

// 更新列的其他信息
func (m *MsSqlDialect) updateColumnOthers(d *sql.DB, table *db.Table) {

	mp := make(map[string]*db.Column, 0)
	for _, v := range table.Columns {
		mp[v.Name] = v
	}

	stmt, err := d.Prepare(fmt.Sprintf(`
	SELECT c.name AS colname,
	ISNULL(p.value,'') AS comment,
	c.is_identity
	FROM sys.columns c
	LEFT JOIN sys.extended_properties p 
	ON p.major_id = c.object_id AND p.minor_id = c.column_id
	WHERE c.object_id=object_id('%s')`, table.Name))
	rows, err := stmt.Query()
	if err != nil {
		log.Println("[ mssql][ dialect]: error ", err.Error())
		return
	}

	for rows.Next() {
		name, comment := "", ""
		isIdentity := false
		err = rows.Scan(&name, &comment, &isIdentity)
		if err == nil {
			if col, ok := mp[name]; ok {
				col.IsAuto = isIdentity
				col.Comment = comment
			}
		} else {
			log.Println("[ mssql][ dialect]: scan error ", err.Error())
		}
	}
	rows.Close()
	stmt.Close()
}

func (m *MsSqlDialect) getTypeId(dbType string) int {
	switch dbType {
	case "smallint":
		return db.TypeInt16
	case "tinyint", "int", "int identity":
		return db.TypeInt32
	case "bit":
		return db.TypeBoolean
	case "bigint", "bigint identity":
		return db.TypeInt64
	case "varbinary", "image", "timestamp":
		return db.TypeBytes
	case "float":
		return db.TypeFloat32
	case "decimal", "numeric", "money", "decimal() identity":
		return db.TypeDecimal
	case "double":
		return db.TypeFloat64
	case "date", "datetime", "smalldatetime":
		return db.TypeDateTime
	case "text", "ntext", "varchar", "nvarchar", "char", "xml", "uniqueidentifier":
		return db.TypeString
	}
	println("[ mssql][ dialect]:Dialect not support type :", dbType)
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
