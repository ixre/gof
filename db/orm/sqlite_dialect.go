/**
 * Copyright 2015 @ at3.net.
 * name : sqlite_dialect
 * author : jarryliu
 * date : 2016-11-11 12:29
 * description :
 * history :
 */
package orm

import (
	"database/sql"
	"fmt"
)

var _ Dialect = new(SqLiteDialect)

type SqLiteDialect struct {
}

func (s *SqLiteDialect) GetField(f string) string {
	return fmt.Sprintf("[%s]", f)
}

func (s *SqLiteDialect) Name() string {
	return "SQLiteDialect"
}

// 获取所有的表
func (s *SqLiteDialect) Tables(db *sql.DB, dbName string, schema string) ([]*Table, error) {
	panic("not implement")
}

// 获取表结构
func (s *SqLiteDialect) Table(db *sql.DB, table string) (*Table, error) {
	panic("not implement")
}
