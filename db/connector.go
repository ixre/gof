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

//create a new connector
func NewSimpleConnector(driverName, driverSource string,
	l log.ILogger, maxConn int, debug bool) Connector {
	db, err := sql.Open(driverName, driverSource)

	if err == nil {
		err = db.Ping()
	}

	if err != nil {
		defer db.Close()
		//如果异常，则显示并退出
		log.Fatalln("[ DBC][ " + driverName + "] " + err.Error())
		return nil
	}

	// 设置最大连接数
	if maxConn > 0 {
		db.SetMaxOpenConns(maxConn)
	}

	return &simpleConnector{
		_db:           db,
		_orm:          orm.NewOrm(db),
		_driverName:   driverName,
		_driverSource: driverName,
		_logger:       l,
	}
}

func (this *simpleConnector) err(err error) error {
	if err != nil {
		if this._logger != nil {
			this._logger.Error(err)
		}
	}
	return err
}

func (this *simpleConnector) debugPrintf(format string, s string, args ...interface{}) {
	if this._debug && this._logger != nil {
		var newArgs []interface{} = make([]interface{}, 0)
		newArgs[0] = s
		newArgs = append(newArgs, args...)
		this._logger.Printf(format+"\n", newArgs...)
	}
}

func (this *simpleConnector) Driver() string {
	return this._driverName
}

func (this *simpleConnector) GetDb() *sql.DB {
	return this._db
}

func (this *simpleConnector) GetOrm() orm.Orm {
	return this._orm
}

func (this *simpleConnector) Query(s string, f func(*sql.Rows), args ...interface{}) error {
	this.debugPrintf("[ DBC][ SQL][ TRACE] - sql = %s ; params = %+v\n", s, args...)
	stmt, err := this.GetDb().Prepare(s)
	if err != nil {
		return this.err(errors.New(fmt.Sprintf(
			"[ DBC][ SQL][ ERROR]:%s ; sql = %s ; params = %+v\n", err.Error(), s, args)))
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		return this.err(errors.New(fmt.Sprintf(
			"[ DBC][ SQL][ ERROR]:%s ; sql = %s ; params = %+v\n", err.Error(), s, args)))
	}
	defer stmt.Close()
	if f != nil {
		f(rows)
	}
	return nil
}

//查询Rows
func (this *simpleConnector) QueryRow(s string, f func(*sql.Row), args ...interface{}) error {
	stmt, err := this.GetDb().Prepare(s)
	if err != nil {
		return this.err(errors.New(fmt.Sprintf(
			"[ DBC][ SQL][ ERROR]:%s ; sql = %s ; params = %+v\n", err.Error(), s, args)))
	} else {
		defer stmt.Close()
		row := stmt.QueryRow(args...)
		if f != nil && row != nil {
			f(row)
		}
	}
	return nil
}

func (this *simpleConnector) ExecScalar(s string, result interface{}, args ...interface{}) (err error) {

	this.debugPrintf("[ DBC][ SQL][ TRACE] - sql = %s ; params = %+v\n", s, args...)

	if result == nil {
		return this.err(errors.New("out result is null pointer."))
	}

	err1 := this.QueryRow(s, func(row *sql.Row) {
		err = row.Scan(result)
	}, args...)

	if err == nil {
		err = err1
	}

	if err != nil {
		return this.err(errors.New(fmt.Sprintf(
			"[ DBC][ SQL][ ERROR]:%s ; sql = %s ; params = %+v\n", err.Error(), s, args)))
	}

	return err
}

//执行
func (this *simpleConnector) Exec(s string, args ...interface{}) (rows int, lastInsertId int, err error) {

	this.debugPrintf("[ DBC][ SQL][ TRACE] - sql = %s ; params = %+v\n", s, args...)

	stmt, err := this.GetDb().Prepare(s)
	if err != nil {
		return 0, -1, err
	}
	result, err := stmt.Exec(args...)
	if err != nil {
		err = this.err(errors.New(fmt.Sprintf(
			"[ DBC][ SQL][ ERROR]:%s ; sql = %s ; params = %+v\n", err.Error(), s, args)))
		return 0, -1, err
	}
	defer stmt.Close()

	lastId, _ := result.LastInsertId()
	affect, _ := result.RowsAffected()

	return int(affect), int(lastId), nil
}

func (this *simpleConnector) ExecNonQuery(sql string, args ...interface{}) (int, error) {
	n, _, err := this.Exec(sql, args...)
	return n, this.err(err)
}
