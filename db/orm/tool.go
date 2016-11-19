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
	"strings"
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

func (t *toolSession) Table2GoStruct(table string) (string, error) {
	tb, err := t.dialect.TableStruct(t.conn, table)
	if err != nil {
		return "", err
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
	return buf.String(), err
}
