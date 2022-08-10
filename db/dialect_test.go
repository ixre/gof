package db

import (
	"testing"
)

func TestPGDialect(t *testing.T) {
	conn, _ := NewConnector("postgresql", "postgres://postgres:123456@127.0.0.1:5432/go2o?sslmode=disable", nil, false)
	tables, err := conn.Dialect().Tables(conn.Raw(), "", "public", "")
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
	conn, _ := NewConnector("mysql", "root:123456@tcp(127.0.0.1:3306)/mysql?charset=utf8", nil, false)
	tables, err := conn.Dialect().Tables(conn.Raw(), "", "", "")
	if err != nil {
		t.Error(err)
	}
	for i := range tables {
		t.Logf("table > %s - %s", tables[i].Name, tables[i].Comment)
	}
}

func TestSQLServerDialect(t *testing.T) {
	conn, err := NewConnector("mssql", "sqlserver://sfDBUser:Jbmeon@008@192.168.16.9:1433?database=DCF19_ERP_TEST_B&encrypt=disable", nil, false)
	if err != nil {
		t.Error(err)
	}
	tables, err := conn.Dialect().Tables(conn.Raw(), "", "", "t_K")
	if err != nil {
		t.Error(err)
	}
	for i := range tables {
		t.Logf("table > %d:%s - %s", i, tables[i].Name, tables[i].Comment)
	}
}
