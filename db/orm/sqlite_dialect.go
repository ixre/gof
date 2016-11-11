/**
 * Copyright 2015 @ at3.net.
 * name : sqlite_dialect
 * author : jarryliu
 * date : 2016-11-11 12:29
 * description :
 * history :
 */
package orm

import "database/sql"

var _ Dialect = new(SqLiteDialect)

type SqLiteDialect struct {
}

func (s *SqLiteDialect) Name() string {
	return "SQLiteDialect"
}

// 获取表结构
func (m *SqLiteDialect) TableStruct(db *sql.DB, table string) (*Table, error) {
	panic("not implement")
}
