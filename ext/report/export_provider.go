package report

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ixre/gof/typeconv"
	"github.com/xuri/excelize/v2"
)

type (
	// ExportProvider 导出
	ExportProvider interface {
		// Export 导出
		Export(rows []map[string]interface{}, fields []string, names []string,
			formatter []ExportFormatter) (binary []byte)
	}
	ExportRow struct {
		Row   map[string]interface{}
		Index int
	}
	// ExportFormatter 数据格式化器
	ExportFormatter func(field string, data interface{}, row *ExportRow) interface{}
)

var (
	// CSV格式
	FCsv = 2
	// Excel格式
	FExcel = 1
	// 文字
	FText = 3
)

// 获取导出提供者
func FactoryExportProvider(format int) ExportProvider {
	switch format {
	case FCsv:
		return &csvProvider{delimer: ","}
	case FText:
		return &textProvider{delimer: ","}
	case FExcel:
		return &excelProvider{}
	}
	panic("not support export format")
}

var _ ExportProvider = new(csvProvider)

type csvProvider struct {
	delimer string
}

func (c *csvProvider) Export(rows []map[string]interface{},
	fields []string, names []string, formatter []ExportFormatter) (binary []byte) {
	buf := bytes.NewBufferString("")
	// 显示表头
	showHeader := len(fields) > 0
	if showHeader {
		for i, k := range names {
			if i > 0 {
				buf.WriteString(c.delimer)
			}
			buf.WriteString(k)
		}
	}
	l := len(rows)
	for i, row := range rows {
		if i < l {
			buf.WriteString("\r\n")
		}
		for fi, field := range fields {
			c.appendField(buf, fi, formatColData(field, row[field], &ExportRow{Row: row, Index: i}))
		}
	}
	return buf.Bytes()
}

func formatColData(field string, data interface{}, row *ExportRow, formatter ...ExportFormatter) interface{} {
	if formatter != nil {
		for _, f := range formatter {
			data = f(field, data, row)
		}
	}
	return data
}

func (c *csvProvider) appendField(buf *bytes.Buffer, ki int, data interface{}) {
	if ki > 0 {
		buf.WriteString(c.delimer)
	}

	dataStr := data.(string)
	specData := strings.Index(dataStr, " ") != -1 ||
		strings.Index(dataStr, "-") != -1 ||
		strings.Index(dataStr, "'") != -1

	if strings.Index(dataStr, "\"") != -1 {
		dataStr = strings.Replace(dataStr, "\"", "\"\"", -1)
		specData = true
	}
	//防止里面含有特殊符号
	if specData {
		buf.WriteString("\"")
		buf.WriteString(dataStr)
		buf.WriteString("\"")
	} else {
		buf.WriteString(dataStr)
	}
}

type textProvider struct {
	delimer string
}

func (t *textProvider) Export(rows []map[string]interface{},
	fields []string, names []string, formatter []ExportFormatter) (binary []byte) {
	buf := bytes.NewBufferString("")
	// 显示表头
	showHeader := fields != nil && len(fields) > 0
	if showHeader {
		for i, k := range names {
			if i > 0 {
				buf.WriteString(t.delimer)
			}
			buf.WriteString(k)
		}
	}
	l := len(rows)
	for i, row := range rows {
		if i < l {
			buf.WriteString("\n")
		}
		for fi, field := range fields {
			t.appendField(buf, fi, formatColData(field, row[field], &ExportRow{Row: row, Index: i}))
		}
	}
	return buf.Bytes()
}

func (t *textProvider) appendField(buf *bytes.Buffer, ki int, data interface{}) {
	if ki > 0 {
		buf.WriteString(t.delimer)
	}
	dataStr := data.(string)
	specData := strings.Index(dataStr, " ") != -1 ||
		strings.Index(dataStr, "-") != -1 ||
		strings.Index(dataStr, "'") != -1

	if strings.Index(dataStr, "\"") != -1 {
		dataStr = strings.Replace(dataStr, "\"", "\"\"", -1)
		specData = true
	}
	//防止里面含有特殊符号
	if specData {
		buf.WriteString("\"")
		buf.WriteString(dataStr)
		buf.WriteString("\"")
	} else {
		buf.WriteString(dataStr)
	}
}

type excelProvider struct {
	csv ExportProvider
}

func NewExcelProvider() ExportProvider {
	return &excelProvider{
		csv: nil,
	}
}

func (e *excelProvider) Export(rows []map[string]interface{},
	fields []string, names []string, formatter []ExportFormatter) (binary []byte) {
	f := excelize.NewFile()
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		log.Println("[ export]: create sheet error ", err)
		return []byte{}
	}
	f.SetActiveSheet(index)
	offset := 0
	// 显示表头
	showHeader := len(fields) > 0
	if showHeader {
		offset = 1
		// 首行表头设置高度为50
		err := f.SetRowHeight("Sheet1", 1, 20)
		if err != nil {
			log.Println("[ export]: set row height error ", err)
		}
		for i, k := range names {
			c := strings.ToUpper(string(rune('A') + int32(i)))
			f.SetCellValue("Sheet1", c+"1", k)
		}

		// 设置宽度
		if len(rows) > 0 {
			for i, k := range fields {
				s := typeconv.Stringify(rows[0][k])
				l := len(s)
				if l < 4 {
					continue
				}
				start := strings.ToUpper(string(rune('A') + int32(i)))
				//end := strings.ToUpper(string(rune('A') + int32(i+1)))
				f.SetColWidth("Sheet1", start, start, float64(l+1)*1.5)
			}
		}
	}
	for i, row := range rows {
		for fi, field := range fields {
			c := strings.ToUpper(string(rune('A') + int32(fi)))
			f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", c, i+offset+1),
				formatColData(field, row[field], &ExportRow{Row: row, Index: i}))
		}
	}
	buf := bytes.NewBuffer(nil)
	f.WriteTo(buf)
	f.Close()
	return buf.Bytes()
}

var Formatter = &formatter{}

type formatter struct {
}

// 格式化日期
func (f *formatter) DateTime(value interface{}) interface{} {
	if s, ok := value.(string); ok {
		return s
	}
	if t, ok := value.(time.Time); ok {
		return t.Format("2006/01/02 15:04")
	}
	if unix, ok := value.(int64); ok {
		return time.Unix(unix, 0).Format("2006/01/02 15:04")
	}
	if unix, ok := value.(int); ok {
		return time.Unix(int64(unix), 0).Format("2006/01/02 15:04")
	}
	return value
}

// 格式化金额
func (f *formatter) IntMoney(value interface{}) interface{} {
	if v, ok := value.(int64); ok {
		return fmt.Sprintf("%.2f", float64(v)/100)
	}
	if v, ok := value.(int); ok {
		return fmt.Sprintf("%.2f", float64(v)/100)
	}
	if v, ok := value.(float64); ok {
		return fmt.Sprintf("%.2f", v/100)
	}
	return value
}

// 占位符
func (f *formatter) Holder(value interface{}, holder string) interface{} {
	if v, ok := value.(string); ok {
		if len(v) == 0 {
			return holder
		}
	}
	return value
}
