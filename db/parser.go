package db

import (
	"database/sql"
	"github.com/ixre/gof/types"
)

// 转换为字典数组
func RowsToMarshalMap(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns() //列名
	//定义数组，数组的类型为[]byte
	var values = make([]interface{}, len(columns))
	for v := range values {
		var a interface{}
		values[v] = &a
	}
	var list []map[string]interface{} //数据切片
	for rows.Next() {
		rows.Scan(values...)
		item := make(map[string]interface{})
		for i, v := range values {
			field := types.CamelTitle(columns[i], true)
			item[field] = *v.(*interface{})
		}
		list = append(list, item)
	}
	return list
}

// 转换为字典数组
// 参考：http://my.oschina.net/nowayout/blog/143278
func RowsToMarshalMapV1(rows *sql.Rows) (rowsMap []map[string]interface{}) {
	rowsMap = []map[string]interface{}{} //数据切片
	var tmpInt = 0                       //序列
	columns, _ := rows.Columns()         //列名

	//定义数组，数组的类型为[]byte
	var values = make([]interface{}, len(columns))
	var rawBytes = make([][]byte, len(values))

	for v := range values {
		values[v] = &rawBytes[v]
	}

	for rows.Next() {
		rows.Scan(values...)

		if len(rowsMap) == tmpInt {
			rowsMap = append(rowsMap, make(map[string]interface{}))
		}

		for i, v := range columns {
			rowsMap[tmpInt][v] = string(rawBytes[i])
		}

		tmpInt++
	}
	return rowsMap
}

func RowsToMap(rows *sql.Rows) (rowsMap []map[string][]byte) {
	rowsMap = []map[string][]byte{} //数据切片
	var tmpInt = 0                  //序列
	columns, _ := rows.Columns()    //列名

	//定义数组，数组的类型为[]byte
	var values = make([]interface{}, len(columns))
	var rawBytes = make([][]byte, len(values))

	for v := range values {
		values[v] = &rawBytes[v]
	}

	for rows.Next() {
		rows.Scan(values...)

		if len(rowsMap) == tmpInt {
			rowsMap = append(rowsMap, make(map[string][]byte))
		}

		for i, v := range columns {
			rowsMap[tmpInt][v] = rawBytes[i]
			//fmt.Println(v + "===>" + string(rawBytes[i]))
		}
		tmpInt++
	}
	return rowsMap
}

func RowToMap(rows *sql.Rows) map[string][]byte {
	rowMap := make(map[string][]byte)
	columns, _ := rows.Columns() //列名
	if rows.Next() {
		row := rows
		//数据
		//定义数组，数组的类型为[]byte
		var values = make([]interface{}, len(columns))
		var rawBytes = make([][]byte, len(values))
		for v := range values {
			values[v] = &rawBytes[v]
		}
		row.Scan(values...)
		for i, v := range columns {
			rowMap[v] = rawBytes[i]
			//fmt.Println(v + "===>" + string(rawBytes[i]))
		}
	}
	return rowMap
}
