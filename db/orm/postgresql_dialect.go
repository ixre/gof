package orm

import (
	"bytes"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

var _ Dialect = new(PostgresqlDialect)

//select datname from pg_database
type PostgresqlDialect struct {
}

func (p *PostgresqlDialect) GetField(f string) string {
	return fmt.Sprintf("\"%s\"", f)
}

func (p *PostgresqlDialect) Name() string {
	return "PostgresqlDialect"
}

// 获取所有的表
func (p *PostgresqlDialect) Tables(db *sql.DB, database string, schema string) ([]*Table, error) {
	//SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'
	buf := bytes.NewBufferString("SELECT table_name FROM information_schema.tables WHERE table_schema ='")
	if schema == "" {
		schema = "public"
	}
	buf.WriteString(schema)
	buf.WriteByte('\'')
	var list []string
	tb := ""
	stmt, err := db.Prepare(buf.String())
	if err == nil {
		rows, err := stmt.Query()
		for rows.Next() {
			if err = rows.Scan(&tb); err == nil {
				list = append(list, tb)
			}
		}
		stmt.Close()
		rows.Close()
		tList := make([]*Table, len(list))
		for i, v := range list {
			if tList[i], err = p.Table(db, v); err != nil {
				return nil, err
			}
			tList[i].Schema = schema
		}
		return tList, nil
	}
	return nil, err
}

// 获取表结构
func (p *PostgresqlDialect) Table(db *sql.DB, table string) (*Table, error) {
	stmt, err := db.Prepare(`SELECT COALESCE(description,'') as comment from pg_description
where objoid='` + table + `'::regclass and objsubid=0`)
	row := stmt.QueryRow()
	comment := ""
	err = row.Scan(&comment)
	if err == nil {
		stmt.Close()
	}
	return p.getStruct(db, table, comment)
}

func (p *PostgresqlDialect) getStruct(db *sql.DB, table, comment string) (*Table, error) {
	//stmt, err := db.Prepare(`SELECT column_name,data_type,udt_name,
	//		is_identity,COALESCE(identity_increment,''),is_nullable
	//		FROM information_schema.columns WHERE table_name ='` + table + `'`)

	smt, err := db.Prepare(strings.Replace(`
SELECT ordinal_position as col_order,column_name,data_type,
coalesce(character_maximum_length,numeric_precision,-1) as col_len,COALESCE(numeric_scale,-1) as col_scale,
CASE is_nullable WHEN 'NO' then 1 else 0 end as not_null,
COALESCE(column_default,'') as col_default,
CASE WHEN position('nextval' in column_default)>0 then 1 else 0 end as is_identity, 
CASE WHEN b.pk_name is null then 0 else 1 end as is_pk,COALESCE(c.DeText,'') as col_comment
FROM information_schema.columns LEFT JOIN (
    SELECT pg_attr.attname as colname,pg_constraint.conname as pk_name from pg_constraint  
    INNER JOIN pg_class on pg_constraint.conrelid = pg_class.oid 
    INNER JOIN pg_attribute pg_attr on pg_attr.attrelid = pg_class.oid and  pg_attr.attnum = pg_constraint.conkey[1] 
    INNER JOIN pg_type on pg_type.oid = pg_attr.atttypid
    WHERE pg_class.relname = '{table}' and pg_constraint.contype='p' 
) b on b.colname = information_schema.columns.column_name
LEFT JOIN (
    select attname,description as DeText from pg_class
    left join pg_attribute pg_attr on pg_attr.attrelid= pg_class.oid
    left join pg_description pg_desc on pg_desc.objoid = pg_attr.attrelid and pg_desc.objsubid=pg_attr.attnum
    where pg_attr.attnum>0 and pg_attr.attrelid=pg_class.oid and pg_class.relname='{table}'
)c on c.attname = information_schema.columns.column_name
where table_schema='public' and table_name='{table}' order by ordinal_position asc
`, "{table}", table, -1))
	var columns []*Column
	colMap := make(map[string]*Column, 0)
	rows, err := smt.Query()
	if err == nil {
		rd := make([]string, 10)
		for rows.Next() {
			if err = rows.Scan(&rd[0], &rd[1], &rd[2], &rd[3], &rd[4], &rd[5], &rd[6], &rd[7], &rd[8], &rd[9]); err == nil {
				len, _ := strconv.Atoi(rd[3])
				c := &Column{
					Name:    strings.TrimSpace(rd[1]),
					IsPk:    rd[8] == "1",
					IsAuto:  rd[7] == "1",
					NotNull: rd[5] == "1",
					DbType:  rd[2],
					Comment: strings.TrimSpace(rd[9]),
					Length:  len,
					Type:    p.getTypeId(rd[2], len),
				}
				//if strings.HasPrefix(table, "wal_") {
				//	println("---", rd[2], len)
				//}
				columns = append(columns, c)
				colMap[c.Name] = c
			}
		}
		smt.Close()
		rows.Close()
		return &Table{
			Name:    table,
			Comment: strings.TrimSpace(comment),
			Engine:  "",
			Charset: "",
			Columns: columns,
		}, err
	}
	return nil, err
}

func (p *PostgresqlDialect) getTypeId(dbType string, len int) int {
	switch dbType {
	case "bigint":
		return TypeInt64
	case "smallint":
		return TypeInt16
	case "numeric", "double precision":
		return TypeFloat64
	case "boolean", "bit":
		return TypeBoolean
	case "text":
		return TypeString
	case "integer":
		if len > 32 {
			return TypeInt64
		} else {
			return TypeInt32
		}
	case "date", "time":
		return TypeDateTime
	}
	if strings.HasPrefix(dbType, "character") {
		return TypeString
	}
	if dbType == "float" {
		if len > 32 {
			return TypeFloat64
		}
		return TypeFloat32
	}
	println("[ ORM][ Postgres][ Warning]:Dialect not support type :", dbType)
	return TypeUnknown
}
