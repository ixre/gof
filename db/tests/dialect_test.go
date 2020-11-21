package tests

import (
	"github.com/ixre/gof/db"
	orm2 "github.com/ixre/gof/db/orm"
	"testing"
)

func TestPGDialect(t *testing.T) {
	conn, _ := db.NewConnector("postgresql", "postgres://postgres:123456@127.0.0.1:5432/go2o?sslmode=disable", nil, false)
	o := orm2.NewOrm(conn.Driver(),conn.Raw())
	tables, err := o.Dialect().Tables(conn.Raw(), "", "public")
	if err != nil {
		t.Error(err)
	}
	for i := range tables {
		tb := tables[i]
		t.Logf("table > %s - %s", tb.Name, tb.Comment)
		for ci := range tb.Columns {
			tc := tb.Columns[ci]
			t.Logf("        -- %s - %s", tc.Name, tc.Comment)
		}
	}
}

func TestMysqlDialect(t *testing.T) {
	conn, _ := db.NewConnector("mysql", "root:123456@tcp(127.0.0.1:3306)/mysql?charset=utf8", nil, false)
	o := orm2.NewOrm(conn.Driver(),conn.Raw())
	tables, err := o.Dialect().Tables(conn.Raw(), "", "")
	if err != nil {
		t.Error(err)
	}
	for i := range tables {
		t.Logf("table > %s - %s", tables[i].Name, tables[i].Comment)
	}
}
