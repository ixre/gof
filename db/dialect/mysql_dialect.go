/**
 * Copyright 2015 @ at3.net.
 * name : mysql_dialect
 * author : jarryliu
 * date : 2016-11-11 12:29
 * description :
 * history :
 */
package dialect

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ixre/gof/db/db"
)

var _ Dialect = new(MySqlDialect)

type MySqlDialect struct {
}

func (m *MySqlDialect) GetField(f string) string {
	if strings.Contains(f, ".") {
		return f
	}
	return fmt.Sprintf("`%s`", f)
}

func (m *MySqlDialect) Name() string {
	return "MySQLDialect"
}

func (m *MySqlDialect) fetchTableNames(d *sql.DB, dbName string, keyword string) ([]string, error) {
	buf := bytes.NewBufferString("SHOW TABLES")
	if dbName != "" {
		buf.WriteString(" FROM `")
		buf.WriteString(dbName)
		buf.WriteString("`")
	}
	if keyword != "" {
		buf.WriteString(` LIKE '%`)
		buf.WriteString(keyword)
		buf.WriteString(`%'`)
	}
	var list []string
	tb := ""
	stmt, err := d.Prepare(buf.String())
	if err == nil {
		if rows, err1 := stmt.Query(); err1 != nil {
			return make([]string, 0), err1
		} else {
			for rows.Next() {
				if err = rows.Scan(&tb); err == nil {
					if strings.HasPrefix(tb, dbName+".") {
						// 当tb包含了库名如：mysql.user会导致可而存在表找不到的情况
						continue
					}
					list = append(list, tb)
				}
			}
			stmt.Close()
			rows.Close()
		}
	}
	return list, err
}

// 获取所有的表
func (m *MySqlDialect) Tables(d *sql.DB, dbName string, keyword string, match func(int, string) bool) (int, []*db.Table, error) {
	tableNames, err := m.fetchTableNames(d, dbName, keyword)
	l := len(tableNames)
	if err != nil {
		return l, nil, err
	}
	tList := make([]*db.Table, 0)
	i := -1
	for _, v := range tableNames {
		i++
		if match != nil && !match(i, v) {
			// 筛选掉不匹配的表
			continue
		}
		t, err2 := m.Table(d, v)
		if err2 != nil {
			return l, nil, err2
		}
		tList = append(tList, t)
	}
	return l, tList, nil

}

// 获取表结构
func (m *MySqlDialect) Table(db *sql.DB, table string) (*db.Table, error) {
	stmt, err := db.Prepare("SHOW CREATE TABLE `" + table + "`")
	if err == nil {
		row := stmt.QueryRow()
		tb, desc := "", ""
		err = row.Scan(&tb, &desc)
		if err == nil {
			stmt.Close()
			return m.getStruct(desc)
		}
	}

	return nil, err
}

func (m *MySqlDialect) getStruct(desc string) (*db.Table, error) {
	/**
	  'mm_member', 'CREATE TABLE `mm_member` (\n  `id` int(11) NOT NULL AUTO_INCREMENT,\n  `user` varchar(20) DEFAULT NULL COMMENT \'用户名\',\n  `pwd` varchar(45) DEFAULT NULL COMMENT \'密码\',\n  `trade_pwd` varchar(45) DEFAULT NULL,\n  `exp` int(11) unsigned DEFAULT \'0\',\n  `level` int(11) DEFAULT \'1\',\n  `invitation_code` varchar(10) DEFAULT NULL COMMENT \'邀请码\',\n  `reg_ip` varchar(20) DEFAULT NULL,\n  `reg_from` varchar(20) DEFAULT NULL,\n  `reg_time` int(11) DEFAULT NULL,\n  `check_code` varchar(8) DEFAULT NULL,\n  `check_expires` int(11) DEFAULT NULL,\n  `login_time` int(11) DEFAULT NULL,\n  `last_login_time` int(11) DEFAULT NULL COMMENT \'最后登录时间\',\n  `state` int(1) DEFAULT \'1\',\n  `update_time` int(11) DEFAULT NULL,\n  PRIMARY KEY (`id`)\n) ENGINE=MyISAM AUTO_INCREMENT=16900 DEFAULT CHARSET=utf8'
	*/
	if desc == "" {
		return nil, errors.New("not found table information")
	}
	//log.Println("-- ErrCode:" + desc+"\n\n")
	//time.Sleep(time.Second)
	i, j := strings.Index(desc, "(\n"), strings.Index(desc, "\n)")
	//获取表名
	tmp := desc[:i]
	name := tmp[strings.Index(tmp, "`")+1 : strings.LastIndex(tmp, "`")]
	//获取表的扩展性信息
	mp := map[string]string{}
	reg := regexp.MustCompile(`\s([^)=]+)=([^\s]+)`)
	matches := reg.FindAllStringSubmatch(desc[j:], -1)
	for _, v := range matches {
		mp[v[1]] = v[2]
	}
	table := &db.Table{
		Name:    name,
		Comment: strings.TrimSpace(mp["COMMENT"]),
		Engine:  mp["ENGINE"],
		Charset: mp["DEFAULT CHARSET"],
		Columns: []*db.Column{},
	}
	if table.Comment != "" {
		table.Comment = strings.Replace(table.Comment, "'", "", -1)
	}
	//获取列信息
	colReg := regexp.MustCompile("`([^`]+)`\\s+([a-z0-9]+[^\\s]+)\\s*")
	commReg := regexp.MustCompile(`.*COMMENT\s*'([^']+)'`)

	colArr := strings.Split(desc[i+3:j], "\n")
	//获取主键
	pkField := ""
	split := "PRIMARY KEY (`"
	i2 := strings.Index(desc, split)
	if i2 != -1 {
		tmp = desc[i2+len(split):]
		pkField = tmp[:strings.Index(tmp, "`")]
	}
	//绑定列
	for _, str := range colArr {
		match := colReg.FindStringSubmatch(str)
		if match != nil {
			dbType := match[2]
			col := &db.Column{
				Name:    match[1],
				DbType:  dbType,
				IsAuto:  strings.Contains(str, "AUTO_"),
				NotNull: strings.Contains(str, "NOT NULL"),
				IsPk:    match[1] == pkField,
				Length:  getTypeLen(dbType),
				Type:    m.getTypeId(dbType),
			}
			comMatch := commReg.FindStringSubmatch(str)
			if comMatch != nil {
				col.Comment = strings.Replace(comMatch[1], "\\n", "", -1)
			}
			table.Columns = append(table.Columns, col)
		}
	}
	return table, nil
}

func (m *MySqlDialect) getTypeId(dbType string) int {
	dbType = strings.ToLower(dbType)
	switch true {
	case strings.HasPrefix(dbType, "smallint"):
		return db.TypeInt16
	case strings.HasPrefix(dbType, "tinyint"):
		return db.TypeInt32
	case strings.HasPrefix(dbType, "bit"):
		return db.TypeBoolean
	case strings.HasPrefix(dbType, "bigint"):
		return db.TypeInt64
	case dbType == "int":
		return db.TypeInt32
	case strings.HasPrefix(dbType, "int("):
		i, _ := strconv.Atoi(dbType[4 : len(dbType)-1])
		if i < 11 {
			return db.TypeInt32
		}
		return db.TypeInt64
	case strings.HasPrefix(dbType, "float"):
		return db.TypeFloat32
	case strings.HasPrefix(dbType, "decimal"):
		return db.TypeDecimal
	case strings.HasPrefix(dbType, "double"):
		return db.TypeFloat64
	case dbType == "datetime", dbType == "date":
		return db.TypeDateTime
	case dbType == "timestamp":
		return db.TypeInt64
	case dbType == "blob", dbType == "longblob":
		return db.TypeBytes
	case dbType == "text", dbType == "longtext",
		dbType == "json", dbType == "mediumtext",
		dbType == "tinytext",
		strings.HasPrefix(dbType, "varchar"),
		strings.HasPrefix(dbType, "char"):
		return db.TypeString
	}
	println("[ ORM][ MySQL][ Warning]:Dialect not support type :", dbType)
	return db.TypeUnknown
}
