/**
 * Copyright 2015 @ at3.net.
 * name : tool.go
 * author : jarryliu
 * date : 2016-11-11 12:19
 * description :
 * history :
 */
package orm

import (
	"bytes"
	"database/sql"
	"github.com/jsix/gof/util"
	"log"
	"strings"
	"text/template"
)

type toolSession struct {
	conn    *sql.DB
	dialect Dialect
}

func NewTool(conn *sql.DB, dialect Dialect) *toolSession {
	return &toolSession{
		conn:    conn,
		dialect: dialect,
	}
}

func (t *toolSession) title(str string) string {
	arr := strings.Split(str, "_")
	for i, v := range arr {
		arr[i] = strings.Title(v)
	}
	return strings.Join(arr, "")
}

func (t *toolSession) goType(dbType string) string {
	l := len(dbType)
	switch true {
	case strings.HasPrefix(dbType, "tinyint"):
		return "int"
	case strings.HasPrefix(dbType, "bit"):
		return "bool"
	case strings.HasPrefix(dbType, "int("):
		if l == 6 {
			return "int"
		}
		return "int64"
	case strings.HasPrefix(dbType, "float"):
		return "float32"
	case strings.HasPrefix(dbType, "decimal"):
		return "float64"
	case dbType == "text", strings.HasPrefix(dbType, "varchar"):
		return "string"
	}
	return "interface{}"
}

// 获取所有的表
func (t *toolSession) Tables(db string) ([]*Table, error) {
	return t.dialect.Tables(t.conn, db)
}

// 获取表结构
func (t *toolSession) Table(table string) (*Table, error) {
	return t.dialect.Table(t.conn, table)
}

// 保存到文件
func (t *toolSession) SaveFile(s string, path string) error {
	return util.BytesToFile([]byte(s), path)
}

// 表生成结构
func (t *toolSession) TableToGoStruct(tb *Table) string {
	if tb == nil {
		return ""
	}
	//log.Println(fmt.Sprintf("%#v", tb))
	buf := bytes.NewBufferString("")
	buf.WriteString("// ")
	buf.WriteString(tb.Comment)
	buf.WriteString("\ntype ")
	buf.WriteString(t.title(tb.Name))
	buf.WriteString(" struct{\n")

	for _, col := range tb.Columns {
		if col.Comment != "" {
			buf.WriteString("    // ")
			buf.WriteString(col.Comment)
			buf.WriteString("\n")
		}
		buf.WriteString("    ")
		buf.WriteString(t.title(col.Name))
		buf.WriteString(" ")
		buf.WriteString(t.goType(col.Type))
		buf.WriteString(" `")
		buf.WriteString("db:\"")
		buf.WriteString(col.Name)
		buf.WriteString("\"")
		if col.Pk {
			buf.WriteString(" pk:\"yes\"")
		}
		if col.Auto {
			buf.WriteString(" auto:\"yes\"")
		}
		buf.WriteString("`")
		buf.WriteString("\n")
	}

	buf.WriteString("}")
	return buf.String()
}

// 表结构生成表单
func (t *toolSession) TableToHtmlForm(tb *Table) string {
	if tb == nil {
		return ""
	}
	//log.Println(fmt.Sprintf("%#v", tb))
	buf := bytes.NewBufferString("")
	buf.WriteString("// ")
	buf.WriteString(tb.Comment)
	buf.WriteString("\ntype ")
	buf.WriteString(t.title(tb.Name))
	buf.WriteString(" struct{\n")

	for _, col := range tb.Columns {
		if col.Comment != "" {
			buf.WriteString("    // ")
			buf.WriteString(col.Comment)
			buf.WriteString("\n")
		}
		buf.WriteString("    ")
		buf.WriteString(t.title(col.Name))
		buf.WriteString(" ")
		buf.WriteString(t.goType(col.Type))
		buf.WriteString(" `")
		buf.WriteString("db:\"")
		buf.WriteString(col.Name)
		buf.WriteString("\"")
		if col.Pk {
			buf.WriteString(" pk:\"yes\"")
		}
		if col.Auto {
			buf.WriteString(" auto:\"yes\"")
		}
		buf.WriteString("`")
		buf.WriteString("\n")
	}

	buf.WriteString("}")
	return buf.String()
}

// 表生成仓储类,sign:函数后是否带签名，ePrefix:实体是否带前缀
func (ts *toolSession) TableToGoRep(tb *Table,
	sign bool, ePrefix string) string {
	if tb == nil {
		return ""
	}

	var err error
	t := &template.Template{}
	t, err = t.Parse(string(TPL_ENTITY_REP))
	if err != nil {
		panic(err)
	}

	pk := "<PK>"
	for i, v := range tb.Columns {
		if i == 0 {
			pk = v.Name
		}
		if v.Pk {
			pk = v.Name
			break
		}
	}
	n := ts.title(tb.Name)
	r2 := ""
	if sign {
		r2 = n
	}
	mp := map[string]interface{}{
		"R":  n + "Rep",
		"R2": r2,
		"E":  n,
		"E2": ePrefix + n,
		"T":  strings.ToLower(tb.Name[:1]),
		"PK": ts.title(pk),
	}
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, mp)
	if err == nil {
		return buf.String()
	}
	log.Println("execute template error:", err.Error())
	return ""
}
