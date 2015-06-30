package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/atnet/gof/db/orm"
	"github.com/atnet/gof/log"
	_ "github.com/go-sql-driver/mysql"
)

var _ Connector = new(SimpleDbConnector)

//数据库连接器
type SimpleDbConnector struct {
	driverName   string  //驱动名称
	driverSource string  //驱动连接地址
	_db          *sql.DB //golang db只需要open一次即可
	_orm         orm.Orm
	logger       log.ILogger
}

//create a new connector
func NewSimpleConnector(driverName, driverSource string,
	l log.ILogger, maxConn int) Connector {
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

	return &SimpleDbConnector{
		_db:          db,
		_orm:         orm.NewOrm(db),
		driverName:   driverName,
		driverSource: driverName,
		logger:       l,
	}
}

func (this *SimpleDbConnector) err(err error) error {
	if err != nil {
		if this.logger != nil {
			this.logger.PrintErr(err)
		}
	}
	return err
}

func (this *SimpleDbConnector) GetDb() *sql.DB {
	return this._db
}

func (this *SimpleDbConnector) GetOrm() orm.Orm {
	return this._orm
}

func (this *SimpleDbConnector) Query(s string, f func(*sql.Rows), args ...interface{}) error {
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
func (this *SimpleDbConnector) QueryRow(s string, f func(*sql.Row), args ...interface{}) error {
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

func (this *SimpleDbConnector) ExecScalar(s string, result interface{}, args ...interface{}) (err error) {
	if result == nil {
		return this.err(errors.New("out result is null pointer."))
	}

	this.QueryRow(s, func(row *sql.Row) {
		err = row.Scan(result)
	}, args...)

	if err != nil {
		return this.err(errors.New(fmt.Sprintf(
			"[ DBC][ SQL][ ERROR]:%s ; sql = %s ; params = %+v\n", err.Error(), s, args)))
	}

	return err
}

//执行
func (this *SimpleDbConnector) Exec(s string, args ...interface{}) (rows int, lastInsertId int, err error) {
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

func (this *SimpleDbConnector) ExecNonQuery(sql string, args ...interface{}) (int, error) {
	n, _, err := this.Exec(sql, args...)
	return n, this.err(err)
}
