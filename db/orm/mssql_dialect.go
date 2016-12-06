/**
 * Copyright 2015 @ at3.net.
 * name : mssql_dialect
 * author : jarryliu
 * date : 2016-11-11 12:29
 * description :
 * history :
 */
package orm

import "database/sql"

var _ Dialect = new(MySqlDialect)

type MsSqlDialect struct {
}

func (m *MsSqlDialect) Name() string {
	return "MSSQLDialect"
}

// 获取所有的表
func (m *MsSqlDialect) Tables(db *sql.DB, dbName string) ([]*Table, error) {
	panic("not implement")
}

// 获取表结构
func (m *MsSqlDialect) Table(db *sql.DB, table string) (*Table, error) {
	panic("not implement")
}
