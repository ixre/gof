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
	TableName     string
	PkFieldName   string
	PkIsAuto      bool
	FieldNames    []string //预留，可能会用到
	FieldMapNames []string
}
