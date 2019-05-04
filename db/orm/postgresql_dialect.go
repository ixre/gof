package orm

import (
	"bytes"
	"database/sql"
)

var _ Dialect = new(PostgresqlDialect)

//select datname from pg_database
type PostgresqlDialect struct {
}

func (p *PostgresqlDialect) Name() string {
	return "PostgresqlDialect"
}

// 获取所有的表
func (p *PostgresqlDialect) Tables(db *sql.DB, database string, schema string) ([]*Table, error) {
	//SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'
	buf := bytes.NewBufferString("SELECT table_name FROM information_schema.tables WHERE table_schema ='")
	if schema != "" {
		buf.WriteString(schema)
	} else {
		buf.WriteString("public")
	}
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
	stmt, err := db.Prepare(`SELECT column_name,data_type,udt_name,
			is_identity,COALESCE(identity_increment,''),is_nullable 
			FROM information_schema.columns WHERE table_name ='` + table + `'`)
	var columns []*Column
	colMap := make(map[string]*Column, 0)
	rows, err := stmt.Query()
	if err == nil {
		rd := make([]string, 6)
		for rows.Next() {
			if err = rows.Scan(&rd[0], &rd[1], &rd[2], &rd[3], &rd[4], &rd[5]); err == nil {
				dbType := rd[2]
				c := &Column{
					Name:    rd[0],
					Pk:      rd[3] == "YES",
					Auto:    rd[4] == "YES",
					NotNull: rd[5] == "YES",
					Type:    rd[2],
					Comment: "",
					Length:  -1,
					GoType:  p.getGoType(rd[1], dbType),
				}
				columns = append(columns, c)
				colMap[c.Name] = c
			}
		}
		stmt.Close()
		rows.Close()
		if stmt, err = db.Prepare(`SELECT b.attname as columnname, COALESCE(a.description,'')  as comment  
 				FROM pg_catalog.pg_description a,pg_catalog.pg_attribute b   
 				WHERE objoid='` + table + `'::regclass AND a.objoid=b.attrelid
				AND a.objsubid=b.attnum`); err == nil {
			rows, err = stmt.Query()
			for rows.Next() {
				if err = rows.Scan(&rd[0], &rd[1]); err == nil {
					if c, ok := colMap[rd[0]]; ok {
						c.Comment = rd[1]
					}
				}
			}
			stmt.Close()
			rows.Close()
		}
	}

	return &Table{
		Name:    table,
		Comment: comment,
		Engine:  "",
		Charset: "",
		Columns: columns,
	}, nil
}

func (p *PostgresqlDialect) getGoType(dbType string, udtName string) int {
	switch udtName {
	case "int2", "int4", "serial", "smallint":
		return GoTypeInt32
	case "boolean", "bit", "bool":
		return GoTypeBoolean
	case "int8", "bigint":
		return GoTypeInt64
	case "float2", "float4":
		return GoTypeFloat32
	case "float8", "money":
		return GoTypeFloat64
	case "varchar":
		return GoTypeString
	}
	return GoTypeUnknown
}
