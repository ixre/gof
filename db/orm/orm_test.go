/**
 * Copyright 2015 @ at3.net.
 * name : orm_test
 * author : jarryliu
 * date : 2016-11-11 15:26
 * description :
 * history :
 */
package orm

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"testing"
)

func getDb() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(172.16.69.128:3306)/mysql?charset=utf8")
	if err == nil {
		err = db.Ping()
	}

	if err != nil {
		defer db.Close()
		//如果异常，则显示并退出
		log.Fatalln("[ DBC][ MySQL] " + err.Error())
		return nil
	}
	return db
}

func TestToolSession_Table2Struct(t *testing.T) {
	d := &MySqlDialect{}
	tool := NewTool(getDb(), d)
	tb, err := tool.Table("user")
	if err != nil {
		t.Error(err)
	}
	str := tool.TableToGoStruct(tb)
	t.Log("//生成的结构代码为：\n" + str + "\n")

	str = tool.TableToGoRepo(tb, true, "model.")
	t.Log("//生成的REP代码为：\n" + str + "\n")
}
