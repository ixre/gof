package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"
)

func TestQuery(t *testing.T) {
	initPGTest()
	println("==== testing query =====")
	values := make([]interface{}, 3)
	scanValues := make([]interface{}, 3)
	for i, v := range values {
		scanValues[i] = &v
	}
	t.Log("testing...")
	_conn.Query("SELECT id,user,pwd FROM mm_member limit 10 offset 0", func(rows *sql.Rows) {
		for rows.Next() {
			rows.Scan(scanValues...)
			s1 := scanValues[0].(*interface{})
			v, ok := (*s1).([]byte)
			t.Log("------", *s1)
			t.Log(v, "=>", ok)
			v2, ok2 := scanValues[1].(string)
			t.Log(v2, "=>", ok2)
		}
	})
}

func TestPGOrmSelect(t *testing.T) {
	initPGTest()
	println("===== testing select model =======")
	for i := 0; i < 3; i++ {
		var users []PGMember
		_orm.Select(&users, "user=$1", "jarry6")
		if i == 0 {
			t.Logf("%#v\n", users)
		}
	}
}

type PGMember struct {
	Id   int32  `db:"id" pk:"yes" auto:"yes"`
	User string `db:"user"`
	Pwd  string `db:"pwd"`
}

func initPGTest() {
	log.SetOutput(os.Stdout)
	_conn,_ = NewConnector("postgresql", "postgres://postgres:123456@127.0.0.1:5432/go2o?sslmode=disable", nil, false)
	_conn.SetMaxIdleConns(0)
	_conn.SetMaxIdleConns(0)
	_conn.SetConnMaxLifetime(time.Second * 10)
}
