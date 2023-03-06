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
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type user struct {
	host string `db:"host"`
}

func getDb() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/mysql?charset=utf8")
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

func TestStmtClose(t *testing.T) {
	db := getDb()
	conn := NewOrm("mysql", db)
	defer db.Close()
	user := []user{}
	err := conn.Select(&user, "1s=?", 1)
	t.Log("---", user, err)
}

//
//func TestToolSession_Table2Struct(t *testing.T) {
//	d := &MySqlDialect{}
//	tool := DialectSession(getDb(), d)
//	tb, err := tool.Table("user")
//	if err != nil {
//		t.Error(err)
//	}
//	str := tool.TableToGoStruct(tb)
//	t.Log("//生成的结构代码为：\n" + str + "\n")
//
//	str = tool.TableToGoRepo(tb, true, "model.")
//	t.Log("//生成的REP代码为：\n" + str + "\n")
//}

func TestGetUnifinedFiles(t *testing.T) {
	s := "select item_info.* FROM item_info"
	o := &simpleOrm{}
	match, txt, prefix := o.getUnifinedField(s)
	t.Log(match, txt, prefix)
}
