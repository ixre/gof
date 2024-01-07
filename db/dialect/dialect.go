package dialect

import (
	"database/sql"
	"regexp"
	"strconv"
	"strings"

	"github.com/ixre/gof/db/db"
)

// // 表筛选条件
// type TableFilterOptions struct {
// 	// 数据库,默认为空
// 	Database string
// 	// 数据库架构
// 	Schema string
// 	// 表名关键词
// 	TableKeyword string
// }

// Dialect 方言
type Dialect interface {
	// 数据库方言名称
	Name() string
	// 获取所有的表
	// filter: 数据表筛选函数，如果为空，则默认匹配
	Tables(db *sql.DB, database string, keyword string, filter func(index int, name string) bool) (int, []*db.Table, error)
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

var lengthRegexp = regexp.MustCompile(`\(([\d|,]+)\)`)

// 获取类型长度
func getTypeLen(dbType string) int {
	// like: l := getTypeLen("varchar(100)")
	// l1 := getTypeLen("decimal(10,2)")
	if lengthRegexp.Match([]byte(dbType)) {
		arr := lengthRegexp.FindAllStringSubmatch(dbType, 1)
		s := strings.Split(arr[0][1], ",")
		i1, err := strconv.Atoi(s[0])
		if err == nil {
			if len(s) == 2 {
				i2, err2 := strconv.Atoi(s[1])
				if err2 != nil {
					panic(err2)
				}
				return i1 + i2
			}
		}
		return i1
	}
	return -1
}
