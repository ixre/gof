package dialect

import (
	"database/sql"

	"github.com/ixre/gof/db/db"
)

// Dialect 方言
type Dialect interface {
	// 数据库方言名称
	Name() string
	// 获取所有的表
	// match: 匹配表名函数，如果为空，则默认匹配
	Tables(db *sql.DB, database string, schema string, match func(string) bool) ([]*db.Table, error)
	// 获取表结构
	Table(db *sql.DB, table string) (*db.Table, error)
	// 获取数据库字段,如果有保留字,则添加引号
	GetField(v string) string
}

// GetDialect 获取方言
func GetDialect(driver string) (string, Dialect) {
	switch driver {
	case "mysql", "mariadb":
		return "mysql", &MySqlDialect{}
	case "postgres", "postgresql", "pgsql":
		return "postgresql", &PostgresqlDialect{}
	case "sqlserver", "mssql":
		return "sqlserver", &MsSqlDialect{}
	case "sqlite":
		return "sqlite", &SqLiteDialect{}
	default:
		panic("不支持的数据库类型" + driver)
	}
}
