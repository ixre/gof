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
)

type dialectSession struct {
	conn    *sql.DB
	dialect Dialect
}

func DialectSession(conn *sql.DB, dialect Dialect) *dialectSession {
	return &dialectSession{
		conn:    conn,
		dialect: dialect,
	}
}

// 获取所有的表
func (d *dialectSession) Tables(db string) ([]*Table, error) {
	return d.dialect.Tables(d.conn, db)
}

// 获取所有的表
func (d *dialectSession) TablesByPrefix(db string,
	prefix string) ([]*Table, error) {
	list, err := d.dialect.Tables(d.conn, db)
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
