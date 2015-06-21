/**
 * Copyright 2013 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-10 21:52
 * description :
 * history :
 */

package orm

type TableMapMeta struct {
	// 表前缀，如果手工添加
	TableName   string
	PkFieldName string
	PkIsAuto    bool
	//字段的索引
	FieldsIndex   []int
	FieldMapNames []string
}
