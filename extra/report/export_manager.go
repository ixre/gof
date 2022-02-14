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
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/ixre/gof/db"
	"log"
	"os"
	"regexp"
	"strings"
)

var _ IDataExportPortal = new(ExportItem)

// ExportItem 导出项目
type ExportItem struct {
	columnMapping []ColumnMapping
	sqlConfig     *ItemConfig
	dbProvider    IDbProvider
	PortalKey     string
}

func (e *ExportItem) formatMappingString(str string) string {
	reg := regexp.MustCompile("\\s|\\n")
	return reg.ReplaceAllString(e.sqlConfig.ColumnMapping, "")
}

// GetColumnMapping 导出的列名(比如：数据表是因为列，这里我需要列出中文列)
func (e *ExportItem) GetColumnMapping() []ColumnMapping {
	if e.columnMapping == nil {
		e.sqlConfig.ColumnMapping = e.formatMappingString(e.sqlConfig.ColumnMapping)
		e.columnMapping = parseColumnMapping(e.sqlConfig.ColumnMapping)
	}
	return e.columnMapping
}

func (e *ExportItem) GetExportColumnNames(
	exportColumns []string) (names []string) {
	names = make([]string, len(exportColumns))
	mapping := e.GetColumnMapping()
	for i, cName := range exportColumns {
		for _, cMap := range mapping {
			if cMap.Field == cName {
				names[i] = cMap.Name
				break
			}
		}
	}
	return names
}

// GetTotalView 获取统计数据
func (e *ExportItem) GetTotalView(ht map[string]string) (row map[string]interface{}) {
	return nil
}

func (e *ExportItem) GetTotalCount(p Params)(int,error) {
	sqlDb := e.dbProvider.GetDB()
	total := 0
	if e.sqlConfig.Total == "" {
		return 0,errors.New("no set total sql")
	}
	sql := SqlFormat(e.sqlConfig.Total, p)
	smt, err := sqlDb.Prepare(e.check(sql))
	if err == nil {
		row := smt.QueryRow()
		smt.Close()
		if row != nil {
			err = row.Scan(&total)
		}
	}
	if err != nil {
		log.Println("[ Export][ Error] -", err.Error(), "\n", sql)
	}
	return total,err
}

func (e *ExportItem) GetSchemaData(p Params)([]map[string]interface{},error){
	if e == nil || e.dbProvider == nil {
		return nil, errors.New("no match config item")
	}
	var sqlRows *sql.Rows
	sqlDb := e.dbProvider.GetDB()
	p.reduce()
	// 获得数据
	if e.sqlConfig.Query == "" {
		return make([]map[string]interface{}, 0),
			errors.New("no such query of item; key:" + e.PortalKey)
	}
	sql := SqlFormat(e.sqlConfig.Query, p)
	//log.Println("-----",sql)
	// 如果包含了多条SQL,那么执行前面SQL语句，查询最后一条语句返回数据
	sqlLines := strings.Split(sql, ";\n")
	if t := len(sqlLines); t > 1 {
		for i, v := range sqlLines {
			if i != t-1 {
				smt, err := sqlDb.Prepare(e.check(v))
				if err == nil {
					smt.Exec()
					smt.Close()
				}
			}
		}
		sql = sqlLines[t-1]
	}
	smt, err := sqlDb.Prepare(e.check(sql))
	if err == nil {
		defer smt.Close()
		sqlRows, err = smt.Query()
		if err == nil {
			data := db.RowsToMarshalMap(sqlRows)
			sqlRows.Close()
			return data,  err
		}
	}
	log.Println("[ Export][ Error] -", err.Error(), "\n", sql)
	return nil,  err
}

func (e *ExportItem) GetSchemaAndData(p Params) ([]map[string]interface{},  int,  error) {
	if e == nil || e.dbProvider == nil {
		return nil, 0, errors.New("no match config item")
	}
	total := -1
	rows,err := e.GetSchemaData(p)
	if err == nil && len(rows) > 0{
		if p.IsFirstIndex(){
			total,err = e.GetTotalCount(p)
		}
	}
	return rows,total,err
}

// GetJsonData 获取要导出的数据Json格式
func (e *ExportItem) GetJsonData(ht map[string]string) string {
	result, err := json.Marshal(nil)
	if err != nil {
		return "{error:'" + err.Error() + "'}"
	}
	return string(result)
}

// Export 导出数据
func (e *ExportItem) Export(parameters *ExportParams,
	provider IExportProvider, formatter IExportFormatter) []byte {
	rows, _, _ := e.GetSchemaAndData(parameters.Params)
	names := e.GetExportColumnNames(parameters.ExportFields)
	fmtArray := []IExportFormatter{interFmt}
	if formatter != nil {
		fmtArray = append(fmtArray, formatter)
	}
	return provider.Export(rows, parameters.ExportFields, names, fmtArray)
}

func (e *ExportItem) check(s string) string {
	if !CheckInject(s) {
		panic("dangers sql: " + s)
	}
	return s
}

// ItemManager 导出项工厂
type ItemManager struct {
	//配置存放路径
	rootPath string
	//配置扩展名
	cfgFileExt string
	//数据库连接
	dbGetter IDbProvider //接口类型不需要加*
	//导出项集合
	exportItems map[string]*ExportItem
	// 是否缓存配置项文件
	cacheFiles bool
}

// NewItemManager 新建导出项目管理器
func NewItemManager(db IDbProvider, rootPath string, cacheFiles bool) *ItemManager {
	if rootPath == "" {
		rootPath = "/query/"
	}
	if rootPath[len(rootPath)-1] != '/' {
		rootPath += "/"
	}
	return &ItemManager{
		dbGetter:    db,
		rootPath:    rootPath,
		cfgFileExt:  ".xml",
		exportItems: make(map[string]*ExportItem),
		cacheFiles:  cacheFiles,
	}
}

// GetItem 获取导出项
func (f *ItemManager) GetItem(portalKey string) IDataExportPortal {
	item, exist := f.exportItems[portalKey]
	if !exist {
		item = f.loadExportItem(portalKey, f.dbGetter)
		if f.cacheFiles {
			f.exportItems[portalKey] = item
		}
	}
	return item
}

// 创建导出项,watch：是否监视文件变化
func (f *ItemManager) loadExportItem(portalKey string,
	db IDbProvider) *ExportItem {
	dir, _ := os.Getwd()
	arr := []string{dir, f.rootPath, portalKey, f.cfgFileExt}
	filePath := strings.Join(arr, "")
	fi, err := os.Stat(filePath)
	if err == nil && fi.IsDir() == false {
		cfg, err1 := readItemConfigFromXml(filePath)
		if err1 == nil {
			return &ExportItem{
				sqlConfig:  cfg,
				PortalKey:  portalKey,
				dbProvider: db,
			}
		}
		err = err1
	}
	if err != nil {
		log.Println("[ Export][ Error]:", err.Error(), "; portal:", portalKey)
	}
	return nil
}

// GetExportData 获取导出数据
func (f *ItemManager) GetExportData(portal string, p Params, page int,
	rows int) (data []map[string]interface{}, total int, err error) {
	exportItem := f.GetItem(portal)
	if exportItem != nil {
		if page > 0 {
			p["page_index"] = page
		}
		if rows > 0 {
			p["page_size"] = rows
		}
		return exportItem.GetSchemaAndData(p)
	}
	return make([]map[string]interface{}, 0), 0, errNoSuchItem
}

// GetWebExportCheckOptions 获取导出列勾选项
func (f *ItemManager) GetWebExportCheckOptions(portal string, token string) (string, error) {
	p := f.GetItem(portal)
	if p == nil {
		return "", errNoSuchItem
	}
	return buildWebExportCheckOptions(p, token), nil
}
