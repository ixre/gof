/**
 * Copyright 2015 @ at3.net.
 * name : tool.go
 * author : jarryliu
 * date : 2016-11-11 12:19
 * description :
 * history :
 */
package orm

import (
	"database/sql"
	"errors"
	"strings"
)

type dialectSession struct {
	conn    *sql.DB
	dialect Dialect
	driver  string
}

func NewDialectSession(db string, conn *sql.DB) (*dialectSession, error) {
	dst := &dialectSession{conn: conn}
	switch db {
	case "mysql", "mariadb":
		dst.driver = "mysql"
		dst.dialect = &MySqlDialect{}
	case "postgres", "postgresql", "pgsql":
		dst.driver = "pgsql"
		dst.dialect = &PostgresqlDialect{}
	case "sqlserver", "mssql":
		dst.driver = "mssql"
		dst.dialect = &MsSqlDialect{}
	case "sqlite":
		dst.driver = "sqlite"
		dst.dialect = &SqLiteDialect{}
	default:
		return nil, errors.New("不支持的数据库类型" + db)
	}
	return dst, nil
}

func DialectSession(conn *sql.DB, dialect Dialect) *dialectSession {
	return &dialectSession{
		conn:    conn,
		dialect: dialect,
	}
}

func (d *dialectSession) Driver() string {
	return d.driver
}

// 获取所有的表
func (d *dialectSession) Tables(database string, schema string) ([]*Table, error) {
	return d.dialect.Tables(d.conn, database, schema)
}

// 获取所有的表
func (d *dialectSession) TablesByPrefix(database string, schema string,
	prefix string) ([]*Table, error) {
	list, err := d.dialect.Tables(d.conn, database, schema)
	if err == nil {
		var l []*Table
		for _, v := range list {
			if strings.HasPrefix(v.Name, prefix) {
				l = append(l, v)
			}
		}
		return l, nil
	}
	return nil, err
}

// 获取表结构
func (d *dialectSession) Table(table string) (*Table, error) {
	return d.dialect.Table(d.conn, table)
}
