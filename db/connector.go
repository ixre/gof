/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-03 20:19
 * description :
 * history :
 */

package db

import (
	"database/sql"
	"github.com/jsix/gof/db/orm"
)

type Connector interface {
	Driver() string

	GetDb() *sql.DB

	GetOrm() orm.Orm

	Query(sql string, f func(*sql.Rows), arg ...interface{}) error

	QueryRow(sql string, f func(*sql.Row), arg ...interface{}) error

	ExecScalar(s string, result interface{}, arg ...interface{}) error

	Exec(sql string, args ...interface{}) (rows int, lastInsertId int, err error)

	ExecNonQuery(sql string, args ...interface{}) (int, error)
}
