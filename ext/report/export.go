/**
 * Copyright 2013 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-02-04 20:13
 * description :
 * history :
 */
package report

import (
	"context"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/ixre/gof/typeconv"
)

var (
	interFmt      = &internalFormatter{}
	errNoSuchItem = errors.New("no such item")
	injectRegexp  = regexp.MustCompile("\\bEXEC\\\\b|UNION.+?SELECT|UPDATE.+?SET|INSERT\\\\s+INTO.+?VALUES|DELETE.+?FROM|(CREATE|ALTER|DROP|TRUNCATE)\\\\s+(TABLE|DATABASE)")
)

type (
	IDbProvider interface {
		// PrepareContext prepare query, clickhouse 应实现QueryRowContext和QueryContext
		//　其他数据库实现PrepareContext即可
		PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
		// QueryRowContext query single row data
		QueryRowContext(todo context.Context, query string) (*sql.Row, error)
		// QueryContext query multiple row data
		QueryContext(todo context.Context, query string) (*sql.Rows, error)
	}

	// ColumnMapping 列映射
	ColumnMapping struct {
		//列的字段
		Field string
		//列的名称
		Name string
	}

	// ItemConfig 导入导出项目配置
	ItemConfig struct {
		ColumnMapping string
		Query         string
		Total         string
		Import        string
	}

	// IDataExportPortal 数据导出入口
	IDataExportPortal interface {
		// GetColumnMapping 导出的列名(比如：数据表是因为列，这里我需要列出中文列)
		GetColumnMapping() []ColumnMapping
		// GetTotalCount 查询总条数
		GetTotalCount(provider IDbProvider, p Params) (int, error)
		// GetSchemaData 查询数据
		GetSchemaData(provider IDbProvider, p Params) ([]map[string]interface{}, error)
		// GetSchemaAndData 获取要导出的数据及表结构,仅在第一页时查询分页数据
		GetSchemaAndData(provider IDbProvider, p Params) (rows []map[string]interface{},
			total int, err error)
		// GetJsonData 获取要导出的数据Json格式
		GetJsonData(ht map[string]string) string
		// GetTotalView 获取统计数据
		GetTotalView(ht map[string]string) (row map[string]interface{})
		// GetExportColumnNames 根据导出的列名获取列的索引及对应键
		GetExportColumnNames(exportColumnNames []string) (fields []string)
		// Export 导出数据
		Export(provider IDbProvider, ep *ExportParams, p IExportProvider, f IExportFormatter) []byte
	}

	// IExportProvider 导出
	IExportProvider interface {
		// Export 导出
		Export(rows []map[string]interface{}, fields []string, names []string,
			formatter []IExportFormatter) (binary []byte)
	}
	// IExportFormatter 数据格式化器
	IExportFormatter interface {
		// Format 格式化字段
		Format(field, name string, rowNum int, data interface{}) interface{}
	}

	// Params 参数
	Params map[string]interface{}

	// ExportParams 导出参数
	ExportParams struct {
		//参数
		Params Params
		//要到导出的列的名称集合
		ExportFields []string
	}
)

const reduceKey = "__reduce"

// Copy 从Map中拷贝数据
func (p Params) Copy(form map[string]string) {
	for k, v := range form {
		if k != "total" && k != "rows" && k != "params" {
			p[k] = strings.TrimSpace(v)
		}
	}
}

// CopyForm 从表单参数中导入数据
func (p Params) CopyForm(form url.Values) {
	for k, v := range form {
		if k != "total" && k != "rows" && k != "params" {
			p[k] = strings.TrimSpace(v[0])
		}
	}
}

func (p Params) IsFirstIndex() bool {
	if !p.Contains(reduceKey) {
		p.reduce()
	}
	return p["page_offset"] == "0"
}

func (p Params) reduce() {
	if p.Contains(reduceKey) {
		return
	}
	//初始化添加参数
	if _, e := p["page_size"]; !e {
		p["page_size"] = "10000000000"
	}
	if _, e := p["page_index"]; !e {
		p["page_index"] = "1"
	}
	// 获取页码和每页加载数量
	pi := p["page_index"]
	ps := p["page_size"]
	pageIndex := typeconv.MustInt(pi)
	pageSize := typeconv.MustInt(ps)
	// 设置SQL分页信息
	if pageIndex > 0 {
		offset := (pageIndex - 1) * pageSize
		p["page_offset"] = strconv.Itoa(offset)
	} else {
		p["page_offset"] = "0"
	}
	p["page_end"] = strconv.Itoa(pageIndex * pageSize)
	p[reduceKey] = true
}

func (p Params) Contains(k string) bool {
	_, ok := p[k]
	return ok
}

// 获取列映射数组
func readItemConfigFromXml(xmlFilePath string) (*ItemConfig, error) {
	var cfg ItemConfig
	content, err := os.ReadFile(xmlFilePath)
	if err != nil {
		return &ItemConfig{}, err
	}
	err = xml.Unmarshal(content, &cfg)
	return &cfg, err
}

// 转换列与字段的映射
func parseColumnMapping(str string) []ColumnMapping {
	re, err := regexp.Compile(`([^:]+):([^;]*);*\s*`)
	if err != nil {
		return nil
	}
	var matches = re.FindAllStringSubmatch(str, -1)
	if matches == nil {
		return nil
	}
	columnsMapping := make([]ColumnMapping, len(matches))
	for i, v := range matches {
		columnsMapping[i] = ColumnMapping{Field: v[1], Name: v[2]}
	}
	return columnsMapping
}

// ParseParams 转换参数
func ParseParams(paramMappings string) Params {
	params := Params{}
	if len(paramMappings) > 0 {
		if paramMappings[0] == '{' {
			if err := json.Unmarshal([]byte(paramMappings), &params); err != nil {
				log.Print("[ export][ param]: parse params failed,"+
					"", err.Error(), "; data=", paramMappings)
			}
		} else {
			paramMappings = strings.Replace(paramMappings,
				"%3d", "=", -1)
			var paramsArr, splitArr []string
			paramsArr = strings.Split(paramMappings, ";")
			//添加传入的参数
			for _, v := range paramsArr {
				splitArr = strings.Split(v, ":")
				params[splitArr[0]] = v[len(splitArr[0])+1:]
			}
		}
	}
	if !params.Contains("where") {
		params["where"] = "0=0"
	}
	return params
}

// CheckInject 判断是否存在危险的注入操作
func CheckInject(s string) bool {
	return !injectRegexp.Match([]byte(s))
}

// SqlFormat 格式化sql语句
func SqlFormat(sql string, ht map[string]interface{}) (formatted string) {
	formatted = sql
	// 替换条件
	//reg := regexp.MustCompile("#if\\s+([^\\s]+)[^\\n]*\n([\\s\\S]+?)#fi\n")
	reg := regexp.MustCompile(`#if\s*[\{|\(](.+?)[\}\)]\s*\n*([\s\S]+?)#fi`)
	submatch := reg.FindAllStringSubmatch(formatted, -1)
	for _, v := range submatch {
		key := v[1]
		dv, ok := ht[key]
		if ok {
			b := checkSqlIf(dv)
			formatted = replaceSqlBlock(formatted, v[0], b, v[2])
			continue
		}
		b := checkIfCompare(ht, key)
		formatted = replaceSqlBlock(formatted, v[0], b, v[2])
	}
	for k, v := range ht {
		formatted = strings.Replace(formatted, "{"+k+"}",
			typeconv.Stringify(v), -1)
	}
	return formatted
}

// 替换SQL中的条件
func replaceSqlBlock(s string, block string, b bool, body string) string {
	i := strings.Index(body, "#else")
	tf, ff := body, ""
	if i != -1 {
		tf = body[:i]
		ff = body[i+6:]
	}
	if b {
		return strings.ReplaceAll(s, block, tf)
	}
	return strings.ReplaceAll(s, block, ff)
}

var mathRegexp = regexp.MustCompile(`(\S+)\s*([><!=]*)\s*(\S+)\s*`)

// 计算判断条件
func checkIfCompare(ht map[string]interface{}, key string) bool {
	submatch := mathRegexp.FindAllStringSubmatch(key, 1)
	for _, match := range submatch {
		vo, ok := ht[match[1]]
		if ok {
			op := match[2]
			if strings.ContainsAny(op, "<>") && op != "<>" {
				return checkIntCompare(op, vo, match[3])
			}
			v1 := typeconv.Stringify(vo)
			v2 := match[3]
			switch op {
			case "=", "==":
				return v1 == v2
			case "<>", "!=":
				return v1 != v2
			}
			//log.Println("----^", match[1], match[2], match[3])
		}
	}
	return false
}

// int条件判断
func checkIntCompare(o string, v1 interface{}, v2 interface{}) bool {
	r := typeconv.MustInt(v1)
	v := typeconv.MustInt(v2)
	switch o {
	case ">":
		return r > v
	case ">=":
		return r >= v
	case "<":
		return r < v
	case "<=":
		return r <= v
	}
	return false
}

// 检查条件是否成立,值为空, false或者小于0均不成立
func checkSqlIf(dv interface{}) bool {
	if dv == "" || dv == "false" || dv == "0" {
		return false
	}
	if dv == false {
		return false
	}
	v, ok := dv.(int)
	if ok && v < 0 {
		return false
	}
	v2, ok2 := dv.(float64)
	if ok2 && v2 < 0 {
		return false
	}
	return true
}

// 内置的格式化器
type internalFormatter struct{}

func (i *internalFormatter) Format(field, name string, rowNum int, data interface{}) interface{} {
	if field == "{row_number}" {
		return strconv.Itoa(rowNum + 1)
	}
	if data == nil {
		return ""
	}
	return data
}
