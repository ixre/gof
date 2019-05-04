package tests

import (
	"github.com/ixre/gof/db"
	"testing"
)

func TestPGDialect(t *testing.T) {
	conn := db.NewConnector("postgresql", "postgres://postgres:123456@127.0.0.1:5432/go2o?sslmode=disable", nil, false)
	orm := conn.GetOrm()
	tables, err := orm.Dialect().Tables(conn.Raw(), "", "public")
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
	conn := db.NewConnector("mysql", "root:123456@tcp(127.0.0.1:3306)/mysql?charset=utf8", nil, false)
	orm := conn.GetOrm()
	tables, err := orm.Dialect().Tables(conn.Raw(), "")
	if err != nil {
		t.Error(err)
	}
	for i := range tables {
		t.Logf("table > %s - %s", tables[i].Name, tables[i].Comment)
	}
}
