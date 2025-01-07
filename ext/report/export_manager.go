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
	"errors"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/ixre/gof/db"
)

var _ IDataExportPortal = new(ExportItem)

// ExportItem 导出项目
type ExportItem struct {
	columnMapping []ColumnMapping
	sqlConfig     *ItemConfig
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

func (e *ExportItem) GetTotalCount(provider IDbProvider, p Params) (int, error) {
	total := 0
	if e.sqlConfig.Total == "" {
		return 0, errors.New("no set total sql")
	}
	query := e.check(SqlFormat(e.sqlConfig.Total, p))
	smt, err := provider.PrepareContext(context.TODO(), query)
	if err == nil {
		var row *sql.Row
		if smt != nil {
			row = smt.QueryRow()
			smt.Close()
		} else {
			// 如果PrepareContext返回stmt为空,error也为空,则直接查询.适用于clickhouse
			row, err = provider.QueryRowContext(context.TODO(), query)
		}
		if row != nil {
			err = row.Scan(&total)
		}
	}
	if err != nil {
		log.Println("[ Export][ Error] -", err.Error(), "\n", query)
	}
	return total, err
}

func (e *ExportItem) GetSchemaData(provider IDbProvider, p Params) ([]map[string]interface{}, error) {
	if e == nil || provider == nil {
		return nil, errors.New("no match config item")
	}
	p.reduce()
	// 获得数据
	if e.sqlConfig.Query == "" {
		return make([]map[string]interface{}, 0),
			errors.New("no such query of item; key:" + e.PortalKey)
	}
	query := SqlFormat(e.sqlConfig.Query, p)
	//log.Println("-----",sql)
	// 如果包含了多条SQL,那么执行前面SQL语句，查询最后一条语句返回数据
	sqlLines := strings.Split(query, ";\n")
	if t := len(sqlLines); t > 1 {
		for i, v := range sqlLines {
			if i != t-1 {
				smt, err := provider.PrepareContext(context.TODO(), e.check(v))
				if err == nil && smt != nil {
					smt.Exec()
					smt.Close()
				}
			}
		}
		query = e.check(sqlLines[t-1])
	}
	smt, err := provider.PrepareContext(context.TODO(), query)
	if err == nil {
		var sqlRows *sql.Rows
		if smt != nil {
			defer smt.Close()
			sqlRows, err = smt.Query()
		} else {
			sqlRows, err = provider.QueryContext(context.TODO(), query)
		}
		if err == nil {
			data := db.RowsToMarshalMap(sqlRows)
			sqlRows.Close()
			return data, err
		}
	}
	log.Println("[ Export][ Error] -", err.Error(), "\n", query)
	return nil, err
}

func (e *ExportItem) GetSchemaAndData(provider IDbProvider, p Params) ([]map[string]interface{}, int, error) {
	if e == nil || provider == nil {
		return nil, 0, errors.New("no match config item")
	}
	total := -1
	rows, err := e.GetSchemaData(provider, p)
	if err == nil && len(rows) > 0 {
		if p.IsFirstIndex() {
			total, err = e.GetTotalCount(provider, p)
		}
	}
	return rows, total, err
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
func (e *ExportItem) Export(db IDbProvider, parameters *ExportParams,
	provider ExportProvider, formatter ExportFormatter) []byte {
	rows, _, _ := e.GetSchemaAndData(db, parameters.Params)
	names := e.GetExportColumnNames(parameters.ExportFields)
	fmtArray := []ExportFormatter{interFmt}
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
	//导出项集合
	exportItems map[string]*ExportItem
	// 是否缓存配置项文件
	cacheFiles bool
	lock       *sync.RWMutex
}

// NewItemManager 新建导出项目管理器
func NewItemManager(rootPath string, cacheFiles bool) *ItemManager {
	if rootPath == "" {
		rootPath = "/query/"
	}
	if rootPath[len(rootPath)-1] != '/' {
		rootPath += "/"
	}
	return &ItemManager{
		rootPath:    rootPath,
		cfgFileExt:  ".xml",
		exportItems: make(map[string]*ExportItem),
		cacheFiles:  cacheFiles,
		lock:        &sync.RWMutex{},
	}
}

// GetItem 获取导出项
func (f *ItemManager) GetItem(portalKey string) IDataExportPortal {
	f.lock.RLock()
	item, exist := f.exportItems[portalKey]
	f.lock.RUnlock()
	if !exist {
		item = f.loadExportItem(portalKey)
		if f.cacheFiles {
			f.lock.Lock()
			f.exportItems[portalKey] = item
			f.lock.Unlock()
		}
	}
	return item
}

// 创建导出项,watch：是否监视文件变化
func (f *ItemManager) loadExportItem(portalKey string) *ExportItem {
	dir, _ := os.Getwd()
	arr := []string{dir, f.rootPath, portalKey, f.cfgFileExt}
	filePath := strings.Join(arr, "")
	fi, err := os.Stat(filePath)
	if err == nil && fi.IsDir() == false {
		cfg, err1 := readItemConfigFromXml(filePath)
		if err1 == nil {
			return &ExportItem{
				sqlConfig: cfg,
				PortalKey: portalKey,
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
func (f *ItemManager) GetExportData(provider IDbProvider, portal string, p Params, page int,
	rows int) (data []map[string]interface{}, total int, err error) {
	exportItem := f.GetItem(portal)
	if exportItem != nil {
		if page > 0 {
			p["page_index"] = page
		}
		if rows > 0 {
			p["page_size"] = rows
		}
		return exportItem.GetSchemaAndData(provider, p)
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
