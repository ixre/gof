package db

import (
	"database/sql"
	"fmt"
	"github.com/jsix/gof/db/orm"
	"log"
	"os"
	"testing"
)

var (
	_connector Connector
	_orm       orm.Orm
	print      bool = false
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
	var usr User
	_orm.Get(&usr, "root")
	println("User:" + usr.User)
	println("Pwd:" + usr.Pwd)
	println("Host:" + usr.Host)
}

func sel() {
	initTest()
	println("===== testing select model =======")
	for i := 0; i < 3; i++ {
		var usrs []User
		_orm.Select(&usrs, "user=?", "root")
		if i == 0 {
			println(usrs)
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
	_connector.Query("SELECT id,usr,pwd FROM mm_member limit 0,10", func(rows *sql.Rows) {
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
//	i, i2, err :=_orm.Save(nil, User{Host: "localhost", User: "uu1", Pwd: "1233455"})
//	fmt.Println(i, i2, err)
//
//	var usr User
//	_orm.Get(&usr, "uu1")
//	fmt.Println("Inserted :", usr)
//
//}

//func Test_savemodel(t *testing.T) {
//	fmt.Println("===== testing save model =======")
//	var usr User
//	_orm.Get(&usr, "uu1")
//	usr.Host = "127.0.0.1"
//	_, _, err := _orm.Save(usr.User, usr)
//	if err != nil {
//		fmt.Println("happend error:", err.Error())
//	} else {
//		_orm.Get(&usr, "uu1")
//		fmt.Println("updated host:", usr.Host)
//	}
//
//}

//func Test_delmodel(t *testing.T) {
//	fmt.Println("===== testing deleting model =======")
//	i, err := _orm.Delete(User{User: "uu1"}, "")
//	fmt.Println(i, "rows deleted")
//	if err != nil {
//		fmt.Println("happend error:", err.Error())
//	}
//}

func initTest() {
	log.SetOutput(os.Stdout)
	_connector = NewSimpleConnector("mysql",
		"root:@tcp(dbs.ts.com:3306)/go2o?charset=utf8", nil, 0, 0, false)
	_orm = _connector.GetOrm()
	_orm.SetTrace(!true)
	_orm.Mapping(User{}, "user")
}
