package db

import (
	"database/sql"
	"fmt"
	"github.com/ixre/gof/db/orm"
	"log"
	"os"
	"testing"
	"time"
)

var (
	_conn Connector
	_orm  orm.Orm
	print = false
)

func repeatRun(fc func(), time int) {
	for i := 0; i < time; i++ {
		fc()
	}
}

func println(args ...interface{}) {
	if print {
		fmt.Println(args...)
	}
}

type User struct {
	User string `db:"user" pk:"yes" auto:"no"`
	Pwd  string `db:"password"`
	Host string `db:"host"`
}

func model() {
	initTest()
	println("===== testing model =======")
	var user User
	_orm.Get(&user, "root")
	println("Username:" + user.User)
	println("Password:" + user.Pwd)
	println("Host:" + user.Host)
}

func sel() {
	initTest()
	println("===== testing select model =======")
	for i := 0; i < 3; i++ {
		var users []User
		_orm.Select(&users, "user=?", "root")
		if i == 0 {
			println(users)
		}
	}
}

func query(t *testing.T) {

	println("==== testing query =====")
	values := make([]interface{}, 3)
	scanValues := make([]interface{}, 3)
	for i, v := range values {
		scanValues[i] = &v
	}
	_conn.Query("SELECT id,user,pwd FROM mm_member limit 0,10", func(rows *sql.Rows) {
		for rows.Next() {
			rows.Scan(scanValues...)

			s1 := scanValues[0].(*interface{})
			v, ok := (*s1).([]byte)
			t.Log("------", *s1)
			t.Log(v, "=>", ok)
			v2, ok2 := scanValues[1].(string)
			t.Log(v2, "=>", ok2)
		}
		//println(RowsToMarshalMap(rows))
	})
}

func Test_to(t *testing.T) {
	initTest()
	repeatRun(func() {
		query(t)
	}, 1)
}

//func Test_model(t *testing.T) {
//	repeatRun(model,10000)
//}

//
//func Test_Select(t *testing.T) {
//	repeatRun(sel,10000)
//}

//
//func Test_insermodel(t *testing.T) {
//
//	fmt.Println("\n===== testing insert model =======")
//	i, i2, err :=_orm.Save(nil, Username{Host: "localhost", Username: "uu1", Password: "1233455"})
//	fmt.Println(i, i2, err)
//
//	var user Username
//	_orm.Get(&user, "uu1")
//	fmt.Println("Inserted :", user)
//
//}

//func Test_savemodel(t *testing.T) {
//	fmt.Println("===== testing save model =======")
//	var user Username
//	_orm.Get(&user, "uu1")
//	user.Host = "127.0.0.1"
//	_, _, err := _orm.Save(user.Username, user)
//	if err != nil {
//		fmt.Println("happend error:", err.Error())
//	} else {
//		_orm.Get(&user, "uu1")
//		fmt.Println("updated host:", user.Host)
//	}
//
//}

//func Test_delmodel(t *testing.T) {
//	fmt.Println("===== testing deleting model =======")
//	i, err := _orm.Delete(Username{Username: "uu1"}, "")
//	fmt.Println(i, "rows deleted")
//	if err != nil {
//		fmt.Println("happend error:", err.Error())
//	}
//}

func initTest() {
	log.SetOutput(os.Stdout)
	_conn, _ = NewConnector("mysql", "root:@tcp(127.0.0.1:3306)/go2o?charset=utf8", nil, false)
	_conn.SetMaxIdleConns(0)
	_conn.SetMaxIdleConns(0)
	_conn.SetConnMaxLifetime(time.Second * 10)
}
