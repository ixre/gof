/**
 * Copyright 2013 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-02-04 20:13
 * description :
 * history :
 */
package report

import (
	_ "database/sql"
	"encoding/xml"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"
)

type (
	//数据项
	DataExportPortal struct {
		//导出的列名(比如：数据表是因为列，这里我需要列出中文列)
		ColumnNames []ColumnMapping
	}

	//列映射
	ColumnMapping struct {
		//列的字段
		Field string
		//列的名称
		Name string
	}

	//导入导出项目配置
	ItemConfig struct {
		ColumnMappingString string
		Query               string
		Total               string
		Import              string
	}

	//数据导出入口
	IDataExportPortal interface {
		//导出的列名(比如：数据表是因为列，这里我需要列出中文列)
		//ColumnNames() (names []DataColumnMapping)
		//获取要导出的数据及表结构
		GetSchemaAndData(ht map[string]string) (rows []map[string]interface{}, total int, err error)
		//获取要导出的数据Json格式
		GetJsonData(ht map[string]string) string
		//获取统计数据
		GetTotalView(ht map[string]string) (row map[string]interface{})
		//根据参数获取导出列名及导出名称
		GetExportColumnIndexAndName(exportColumnNames []string) (dict map[string]string)
	}

	//导出
	IDataExportProvider interface {
		//导出
		Export(rows []map[string]interface{}, columns map[string]string) (binary []byte)
	}

	//导出参数
	Params struct {
		//参数
		Parameters map[string]string
		//要到导出的列(对应IDataExportPortal的ColumnNames或DataTable的Shema
		ExportColumnNames []string
	}
)

// 从Map中拷贝数据
func (this *Params) Copy(form map[string]string) {
	for k, v := range form {
		if k != "total" && k != "rows" && k != "params" {
			this.Parameters[k] = v
		}
	}
}

// 从表单参数中导入数据
func (this *Params) CopyForm(form url.Values) {
	for k, v := range form {
		if k != "total" && k != "rows" && k != "params" {
			this.Parameters[k] = v[0]
		}
	}
}

//根据参数获取导出列名及导出名称
func (portal *DataExportPortal) GetExportColumnIndexAndName(
	exportColumnNames []string) (dict map[string]string) {
	dict = make(map[string]string)
	for _, cName := range exportColumnNames {
		for _, cMap := range portal.ColumnNames {
			if cMap.Name == cName {
				dict[cMap.Field] = cMap.Name
				break
			}
		}
	}
	return dict
}

//获取列映射
func GetColumnMappings(columnMappingString string) (
	columnsMapping []ColumnMapping, err error) {
	re, err := regexp.Compile("([^:]+):([;]*)")
	if err != nil {
		return nil, err
	}

	var matches [][]string = re.FindAllStringSubmatch(columnMappingString, 0)
	if matches == nil {
		return nil, nil
	}
	columnsMapping = make([]ColumnMapping, 0, len(matches))
	for i, v := range matches {
		columnsMapping[i] = ColumnMapping{Field: v[1], Name: v[2]}
	}
	return columnsMapping, nil
}

//获取列映射数组
func LoadExportConfigFromXml(xmlFilePath string) (*ItemConfig, error) {
	var cfg ItemConfig
	content, _err := ioutil.ReadFile(xmlFilePath)
	if _err != nil {
		return &ItemConfig{}, _err
	}
	err := xml.Unmarshal(content, &cfg)
	return &cfg, err
}

func Export(portal IDataExportPortal, parameters Params,
	provider IDataExportProvider) []byte {
	rows, _, _ := portal.GetSchemaAndData(parameters.Parameters)
	return provider.Export(rows, portal.GetExportColumnIndexAndName(
		parameters.ExportColumnNames))
}

func GetExportParams(paramMappings string, columnNames []string) *Params {

	var parameters map[string]string = make(map[string]string)

	if paramMappings != "" {

		paramMappings = strings.Replace(paramMappings, "%3d", "=", -1)
		var paramsArr, splitArr []string

		paramsArr = strings.Split(paramMappings, ";")

		//添加传入的参数
		for _, v := range paramsArr {
			splitArr = strings.Split(v, ":")
			parameters[splitArr[0]] = v[len(splitArr[0])+1:]
		}

	}
	return &Params{ExportColumnNames: columnNames, Parameters: parameters}

}

// 格式化sql语句
func SqlFormat(sql string, ht map[string]string) (formatted string) {
	formatted = sql
	for k, v := range ht {
		formatted = strings.Replace(formatted, "{"+k+"}", v, -1)
	}
	return formatted
}
