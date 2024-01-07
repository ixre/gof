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
	"strings"

	"github.com/ixre/gof/db/db"
	"github.com/ixre/gof/db/dialect"
)

type dialectSession struct {
	conn    *sql.DB
	dialect dialect.Dialect
	driver  string
}

func NewDialectSession(db string, conn *sql.DB) (*dialectSession, error) {
	dst := &dialectSession{conn: conn}
	dst.driver, dst.dialect = dialect.GetDialect(db)
	return dst, nil
}

func DialectSession(conn *sql.DB, dialect dialect.Dialect) *dialectSession {
	return &dialectSession{
		conn:    conn,
		dialect: dialect,
	}
}

func (d *dialectSession) Driver() string {
	return d.driver
}

// 获取所有的表
func (d *dialectSession) Tables(database string, keyword string, match func(int, string) bool) (int, []*db.Table, error) {
	return d.dialect.Tables(d.conn, database, keyword, match)
}

// 获取所有的表
func (d *dialectSession) TablesByPrefix(database string, schema string,
	prefix string) ([]*db.Table, error) {
	_, tables, err := d.dialect.Tables(d.conn, database, schema, func(i int, s string) bool {
		return strings.HasPrefix(s, prefix)
	})
	return tables, err
}

// 获取表结构
func (d *dialectSession) Table(table string) (*db.Table, error) {
	return d.dialect.Table(d.conn, table)
}
