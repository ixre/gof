package orm

import (
	"bytes"
	"database/sql"
)

var _ Dialect = new(PostgresqlDialect)

//select datname from pg_database
type PostgresqlDialect struct {
}

func (m *PostgresqlDialect) Name() string {
	return "PostgresqlDialect"
}

// 获取所有的表
func (m *PostgresqlDialect) Tables(db *sql.DB, schemaName string) ([]*Table, error) {
	//SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'
	buf := bytes.NewBufferString("SELECT table_name FROM information_schema.tables WHERE table_schema ='")
	if schemaName != "" {
		buf.WriteString(schemaName)
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
			if tList[i], err = m.Table(db, v); err != nil {
				return nil, err
			}
		}
		return tList, nil
	}
	return nil, err
}

// 获取表结构
func (m *PostgresqlDialect) Table(db *sql.DB, table string) (*Table, error) {
	stmt, err := db.Prepare(`SELECT COALESCE(description,'') as comment from pg_description
where objoid='` + table + `'::regclass and objsubid=0`)
	row := stmt.QueryRow()
	comment := ""
	err = row.Scan(&comment)
	if err == nil {
		stmt.Close()
	}
	return m.getStruct(db, table, comment)
}

func (m *PostgresqlDialect) getStruct(db *sql.DB, table, comment string) (*Table, error) {
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
				c := &Column{
					Name:    rd[0],
					Pk:      rd[3] == "YES",
					Auto:    rd[4] == "YES",
					NotNull: rd[5] == "YES",
					Type:    rd[2],
					Comment: "",
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
