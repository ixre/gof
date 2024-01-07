package db

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestPGDialect(t *testing.T) {
	conn, _ := NewConnector("postgresql", "postgres://postgres:123456@127.0.0.1:5432/go2o?sslmode=disable", nil, false)
	_, tables, err := conn.Dialect().Tables(conn.Raw(), "", "public", nil)
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
	conn, err1 := NewConnector("mysql", "mzl:mzl888@tcp(127.0.0.1:3306)/mzl-next?charset=utf8", nil, false)
	if err1 != nil {
		t.Error(err1)
		t.FailNow()
	}
	match := func(i int, s string) bool {
		return i > 2 && i <= 7
	}
	_, tables, err := conn.Dialect().Tables(conn.Raw(), "", "uams_", match)
	if err != nil {
		t.Error(err)
	}
	for i := range tables {
		t.Logf("table > %s - %s", tables[i].Name, tables[i].Comment)
	}
}

func TestSQLServerDialect(t *testing.T) {
	conn, err := NewConnector("mssql", "sqlserver://sfDBUser:Jbmeon@008@192.168.16.119:1433?database=DCF19_ERP_TEST_B&encrypt=disable", nil, false)
	if err != nil {
		t.Error(err)
	}
	_, tables, err := conn.Dialect().Tables(conn.Raw(), "", "", func(i int, s string) bool { return s == "t_COPD_OrdMst" })
	if err != nil {
		t.Error(err)
	}
	for i := range tables {
		t.Logf("table > %d:%s - %s", i, tables[i].Name, tables[i].Comment)
	}
}
