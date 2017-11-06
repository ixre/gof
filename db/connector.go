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
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jsix/gof/db/orm"
	"github.com/jsix/gof/log"
)

type (
	Connector interface {
		Driver() string

		GetDb() *sql.DB

		GetOrm() orm.Orm

		Query(sql string, f func(*sql.Rows), arg ...interface{}) error

		QueryRow(sql string, f func(*sql.Row) error, arg ...interface{}) error

		ExecScalar(s string, result interface{}, arg ...interface{}) error

		Exec(sql string, args ...interface{}) (rows int, lastInsertId int, err error)

		ExecNonQuery(sql string, args ...interface{}) (int, error)
	}
)

var _ Connector = new(simpleConnector)

//数据库连接器
type simpleConnector struct {
	_driverName   string  //驱动名称
	_driverSource string  //驱动连接地址
	_db           *sql.DB //golang db只需要open一次即可
	_orm          orm.Orm
	_logger       log.ILogger
	_debug        bool // 是否调试模式
}

// create a new connector
func NewSimpleConnector(driverName, driverSource string,
	l log.ILogger, maxConn int, debug bool) Connector {
	db, err := sql.Open(driverName, driverSource)
	if err == nil {
		err = db.Ping()
	}
	if err != nil {
		db.Close()
		//如果异常，则显示并退出
		log.Fatalln("[ DBC][ " + driverName + "] " + err.Error())
		return nil
	}

	// 设置最大连接数,设置MaxOpenConns和MaxIdleConns
	// 不出现：statement.go:27: Invalid Connection 警告信息
	if maxConn > 0 {
		db.SetMaxOpenConns(maxConn)
		db.SetMaxIdleConns(maxConn)
	}

	o := orm.NewOrm(driverName, db)
	if debug {
		o.SetTrace(true)
	}
	return &simpleConnector{
		_db:           db,
		_orm:          o,
		_driverName:   driverName,
		_driverSource: driverName,
		_logger:       l,
		_debug:        debug,
	}
}

func (t *simpleConnector) err(err error) error {
	if err != nil {
		if t._logger != nil {
			t._logger.Error(err)
		}
	}
	return err
}

func (t *simpleConnector) debugPrintf(format string, s string, args ...interface{}) {
	if t._debug && t._logger != nil {
		newArgs := []interface{}{s}
		newArgs = append(newArgs, args...)
		t._logger.Printf(format+"\n", newArgs...)
	}
}

func (t *simpleConnector) Driver() string {
	return t._driverName
}

func (t *simpleConnector) GetDb() *sql.DB {
	return t._db
}

func (t *simpleConnector) GetOrm() orm.Orm {
	return t._orm
}

func (t *simpleConnector) Query(s string, f func(*sql.Rows), args ...interface{}) error {
	t.debugPrintf("[ DBC][ SQL][ TRACE] - sql = %s ; params = %+v\n", s, args)
	stmt, err := t.GetDb().Prepare(s)
	var rows *sql.Rows
	if err == nil {
		rows, err = stmt.Query(args...)
	}
	if err == nil {
		stmt.Close()
		defer rows.Close()
		if f != nil && rows != nil {
			f(rows)
		}
	} else if err != sql.ErrNoRows {
		err = t.err(errors.New(fmt.Sprintf(
			"[ DBC][ SQL][ ERROR]:%s ; sql = %s ; params = %+v\n",
			err.Error(), s, args)))
	}
	return err
}

//查询Rows
func (t *simpleConnector) QueryRow(s string, f func(*sql.Row) error, args ...interface{}) error {
	stmt, err := t.GetDb().Prepare(s)
	if err != nil {
		return t.err(errors.New(fmt.Sprintf(
			"[ DBC][ SQL][ ERROR]:%s ; sql = %s ; params = %+v\n", err.Error(), s, args)))
	} else {
		defer stmt.Close()
		row := stmt.QueryRow(args...)
		if f != nil && row != nil {
			return f(row)
		}
	}
	return err
}

func (t *simpleConnector) ExecScalar(s string, result interface{},
	args ...interface{}) (err error) {
	t.debugPrintf("[ DBC][ SQL][ TRACE] - sql = %s ; params = %+v\n", s, args)
	if result == nil {
		return t.err(errors.New("out result is null"))
	}
	err = t.QueryRow(s, func(row *sql.Row) error {
		return row.Scan(result)
	}, args...)
	if err != nil && err != sql.ErrNoRows {
		return t.err(errors.New(fmt.Sprintf(
			"[ DBC][ SQL][ ERROR]:%s ; sql = %s ; params = %+v\n", err.Error(), s, args)))
	}
	return err
}

//执行
func (t *simpleConnector) Exec(s string, args ...interface{}) (rows int, lastInsertId int, err error) {
	t.debugPrintf("[ DBC][ SQL][ TRACE] - sql = %s ; params = %+v\n", s, args)
	stmt, err := t.GetDb().Prepare(s)
	if err != nil {
		return 0, -1, err
	}
	result, err := stmt.Exec(args...)
	if err != nil {
		err = t.err(errors.New(fmt.Sprintf(
			"[ DBC][ SQL][ ERROR]:%s ; sql = %s ; params = %+v\n", err.Error(), s, args)))
		return 0, -1, err
	}
	var lastId int64
	var affect int64
	affect, err = result.RowsAffected()
	if err == nil {
		stmt.Close()
		lastId, err = result.LastInsertId()
	}
	return int(affect), int(lastId), err
}

func (t *simpleConnector) ExecNonQuery(sql string, args ...interface{}) (int, error) {
	n, _, err := t.Exec(sql, args...)
	return n, t.err(err)
}
