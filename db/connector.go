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
	"time"
)

type (
	Connector interface {
		Driver() string
		Raw() *sql.DB
		GetOrm() orm.Orm
		SetMaxOpenConns(n int)
		SetMaxIdleConns(n int)
		SetConnMaxLifetime(d time.Duration)
		Query(sql string, f func(*sql.Rows), arg ...interface{}) error
		QueryRow(sql string, f func(*sql.Row) error, arg ...interface{}) error
		ExecScalar(s string, result interface{}, arg ...interface{}) error
		Exec(sql string, args ...interface{}) (rows int, lastInsertId int, err error)
		ExecNonQuery(sql string, args ...interface{}) (int, error)
	}
)

var _ Connector = new(defaultConnector)

// create a new connector
func NewConnector(driverName, driverSource string, l log.ILogger, debug bool) Connector {
	db, err := sql.Open(driverName, driverSource)
	if err == nil {
		err = db.Ping()
	}
	if err != nil {
		db.Close()
		//如果异常，则显示并退出
		log.Fatalln("[ Gof][ Connector]:" + driverName + "-" + err.Error())
		return nil
	}
	o := orm.NewOrm(driverName, db)
	if debug {
		o.SetTrace(true)
	}
	return &defaultConnector{
		db:           db,
		orm:          o,
		driverName:   driverName,
		driverSource: driverName,
		logger:       l,
		debug:        debug,
	}
}

//数据库连接器
type defaultConnector struct {
	driverName   string  //驱动名称
	driverSource string  //驱动连接地址
	db           *sql.DB //golang db只需要open一次即可
	orm          orm.Orm
	logger       log.ILogger
	debug        bool // 是否调试模式
}

func (t *defaultConnector) err(err error) error {
	if err != nil {
		if t.logger != nil {
			t.logger.Error(err)
		}
	}
	return err
}

func (t *defaultConnector) debugPrintf(format string, s string, args ...interface{}) {
	if t.debug && t.logger != nil {
		newArgs := []interface{}{s}
		newArgs = append(newArgs, args...)
		t.logger.Printf(format+"\n", newArgs...)
	}
}

func (t *defaultConnector) Driver() string {
	return t.driverName
}

func (t *defaultConnector) Raw() *sql.DB {
	return t.db
}

func (t *defaultConnector) GetOrm() orm.Orm {
	return t.orm
}

// 设置最大打开的连接数
func (t *defaultConnector) SetMaxOpenConns(n int) {
	t.db.SetMaxOpenConns(n)
}

// 设置最大闲置的连接数
func (t *defaultConnector) SetMaxIdleConns(n int) {
	t.db.SetMaxIdleConns(n)
}

// 设置连接存活时间,同Mysql的wait_timeout
func (t *defaultConnector) SetConnMaxLifetime(d time.Duration) {
	t.db.SetConnMaxLifetime(d)
}

func (t *defaultConnector) Query(s string, f func(*sql.Rows), args ...interface{}) error {
	t.debugPrintf("[ DBC][ SQL][ TRACE] - sql = %s ; params = %+v\n", s, args)
	stmt, err := t.Raw().Prepare(s)
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
func (t *defaultConnector) QueryRow(s string, f func(*sql.Row) error, args ...interface{}) error {
	stmt, err := t.Raw().Prepare(s)
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

func (t *defaultConnector) ExecScalar(s string, result interface{},
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
func (t *defaultConnector) Exec(s string, args ...interface{}) (rows int, lastInsertId int, err error) {
	t.debugPrintf("[ DBC][ SQL][ TRACE] - sql = %s ; params = %+v\n", s, args)
	stmt, err := t.Raw().Prepare(s)
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

func (t *defaultConnector) ExecNonQuery(sql string, args ...interface{}) (int, error) {
	n, _, err := t.Exec(sql, args...)
	return n, t.err(err)
}
