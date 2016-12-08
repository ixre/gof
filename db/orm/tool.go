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
	"regexp"
	"strings"
	"text/template"
)

var (
	emptyReg       = regexp.MustCompile("\\s+\"\\s*\"\\s*\\n")
	emptyImportReg = regexp.MustCompile("import\\s*\\(([\\n\\s\"]+)\\)")
)

const (
	//模型包名
	V_ModelPkgName = "ModelPkgName"
	//仓储结构包名
	V_RepoPkgName = "RepoPkgName"
	//仓储接口包名
	V_IRepoPkgName = "IRepoPkgName"
	//仓储结构引用模型包路径
	V_ModelPkg = "ModelPkg"
	//仓储接口引用模型包路径
	V_ModelPkgIRepo = "ModelPkgIRepo"
)

type toolSession struct {
	conn    *sql.DB
	dialect Dialect
	//生成代码变量
	codeVars map[string]interface{}
}

func NewTool(conn *sql.DB, dialect Dialect) *toolSession {
	return (&toolSession{
		conn:     conn,
		dialect:  dialect,
		codeVars: make(map[string]interface{}),
	}).init()
}

func (t *toolSession) init() *toolSession {
	t.Var(V_ModelPkgName, "model")
	t.Var(V_RepoPkgName, "repo")
	t.Var(V_IRepoPkgName, "repo")
	t.Var(V_ModelPkg, "")
	t.Var(V_ModelPkgIRepo, "")
	return t
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

// 获取所有的表
func (t *toolSession) TablesByPrefix(db string, prefix string) ([]*Table, error) {
	list, err := t.dialect.Tables(t.conn, db)
	if err == nil {
		l := []*Table{}
		for _, v := range list {
			if strings.HasPrefix(v.Name, prefix) {
				l = append(l, v)
			}
		}
		return l, nil
	}
	return nil, err
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
	pkgName := ""
	if p, ok := t.codeVars[V_ModelPkgName]; ok {
		pkgName = p.(string)
	} else {
		pkgName = "model"
	}

	//log.Println(fmt.Sprintf("%#v", tb))
	buf := bytes.NewBufferString("")
	buf.WriteString("package ")
	buf.WriteString(pkgName)

	buf.WriteString("\n// ")
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

// 解析模板
func (ts *toolSession) Resolve(t CodeTemplate) CodeTemplate {
	t = resolveRepTag(t)
	return t
}

// 定义变量或修改变量
func (ts *toolSession) Var(key string, v interface{}) {
	if v == nil {
		delete(ts.codeVars, key)
		return
	}
	ts.codeVars[key] = v
}

// 生成代码
func (ts *toolSession) generateCode(tb *Table, tpl CodeTemplate,
	sign bool, ePrefix string) string {
	if tb == nil {
		return ""
	}

	var err error
	t := &template.Template{}
	t, err = t.Parse(string(tpl))
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
		"VAR": ts.codeVars,
		"R":   n + "Rep",
		"R2":  r2,
		"E":   n,
		"E2":  ePrefix + n,
		"T":   strings.ToLower(tb.Name[:1]),
		"PK":  ts.title(pk),
	}
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, mp)
	if err == nil {
		code := buf.String()
		//去除空引用
		code = emptyImportReg.ReplaceAllString(code, "")
		//如果不包含模型，则可能为引用空的包
		code = emptyReg.ReplaceAllString(code, "")
		return code
	}
	log.Println("execute template error:", err.Error())
	return ""
}

// 表生成仓储结构,sign:函数后是否带签名，ePrefix:实体是否带前缀
func (ts *toolSession) TableToGoRepo(tb *Table,
	sign bool, ePrefix string) string {
	return ts.generateCode(tb, TPL_ENTITY_REP, sign, ePrefix)
}

// 表生成仓库仓储接口
func (ts *toolSession) TableToGoIRepo(tb *Table,
	sign bool, ePrefix string) string {
	return ts.generateCode(tb, TPL_ENTITY_REP_INTERFACE, sign, ePrefix)
}
