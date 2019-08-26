/**
 * Copyright 2015 @ at3.net.
 * name : mysql_dialect
 * author : jarryliu
 * date : 2016-11-11 12:29
 * description :
 * history :
 */
package orm

import (
	"bytes"
	"database/sql"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var _ Dialect = new(MySqlDialect)

type MySqlDialect struct {
}

func (m *MySqlDialect) Name() string {
	return "MySQLDialect"
}

// 获取所有的表
func (m *MySqlDialect) Tables(db *sql.DB, dbName string, schema string) ([]*Table, error) {
	buf := bytes.NewBufferString("SHOW TABLES")
	if dbName != "" {
		buf.WriteString(" FROM `")
		buf.WriteString(dbName)
		buf.WriteString("`;")
	}
	var list []string
	tb := ""
	stmt, err := db.Prepare(buf.String())
	if err == nil {
		rows, err := stmt.Query()
		for rows.Next() {
			if err = rows.Scan(&tb); err == nil {
				list = append(list, tb)
			}
		}
		stmt.Close()
		rows.Close()
		tList := make([]*Table, len(list))
		for i, v := range list {
			if tList[i], err = m.Table(db, v); err != nil {
				return nil, err
			}
		}
		return tList, nil
	}
	return nil, err
}

// 获取表结构
func (m *MySqlDialect) Table(db *sql.DB, table string) (*Table, error) {
	stmt, err := db.Prepare("SHOW CREATE TABLE " + table)
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

func (m *MySqlDialect) getStruct(desc string) (*Table, error) {
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
	reg := regexp.MustCompile("\\s([^)=]+)=([^\\s]+)")
	matches := reg.FindAllStringSubmatch(desc[j:], -1)
	for _, v := range matches {
		mp[v[1]] = v[2]
	}
	table := &Table{
		Name:    name,
		Comment: mp["COMMENT"],
		Engine:  mp["ENGINE"],
		Charset: mp["DEFAULT CHARSET"],
		Columns: []*Column{},
	}
	if table.Comment != "" {
		table.Comment = strings.Replace(table.Comment, "'", "", -1)
	}

	//获取列信息
	colReg := regexp.MustCompile("`([^`]+)`\\s+([a-z0-9]+[^\\s]+)\\s")
	commReg := regexp.MustCompile("COMMENT\\s\\\\*'([^']+)'")
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
			col := &Column{
				Name:    match[1],
				Type:    dbType,
				IsAuto:    strings.Index(str, "AUTO_") != -1,
				NotNull: strings.Index(str, "NOT NULL") != -1,
				IsPk:    match[1] == pkField,
				Length:  m.getTypeLen(dbType),
				TypeId:  m.getTypeId(dbType),
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
	switch true {
	case strings.HasPrefix(dbType, "tinyint"):
		return TypeInt32
	case strings.HasPrefix(dbType, "bit"):
		return TypeBoolean
	case dbType == "bigint":
		return TypeInt64
	case strings.HasPrefix(dbType, "int("):
		//i, _ := strconv.Atoi(dbType[4 : len(dbType)-1])
		//if i <= 11 {
		return TypeInt32
		//}
		//return TypeInt64
	case strings.HasPrefix(dbType, "float"):
		return TypeFloat32
	case strings.HasPrefix(dbType, "decimal"),
		strings.HasPrefix(dbType, "double"):
		return TypeFloat64
	case dbType == "text", dbType == "longtext",
		strings.HasPrefix(dbType, "varchar"):
		return TypeString
	}
	println("[ ORM][ MySQL][ Warning]:Dialect not support type :", dbType)
	return TypeUnknown
}

// 获取类型长度
func (m *MySqlDialect) getTypeLen(dbType string) int {
	i := strings.Index(dbType, "(")
	j := strings.LastIndex(dbType, ")")
	if i > 0 && j > 0 {
		s := strings.Split(dbType[i+1:j], ",")
		i1, err := strconv.Atoi(s[0])
		if err != nil {
			panic(err)
		}
		if len(s) == 2 {
			i2, err2 := strconv.Atoi(s[1])
			if err2 != nil {
				panic(err2)
			}
			return i1 + i2
		}
		return i1
	}
	return -1
}
