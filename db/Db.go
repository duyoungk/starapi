package db

import (
	"container/list"
	"database/sql"
	"fmt"
	"starapi/util"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
)

type Query struct {
	dbId       int
	connString string
	db         *sql.DB
}

var UserDb []*Query
var MainDB *Query

func init() {
	fc := &util.FileConfig{}
	fc.LoadFile("dbconn.ini")

	MainDB = &Query{}
	mainConnStr := fc.Get("conn_main", "")
	err := MainDB.Open(0, mainConnStr.(string))
	if err != nil {
		fmt.Println("CANNOT OPEN MAIN DB:", err.Error())
	}

	connstrs := new(list.List)

	for i := 1; i <= 100; i++ {
		key := fmt.Sprintf("conn_user%03d", i)
		connstr := fc.Get(key, "")
		if connstr == nil || connstr == "" {
			break
		}
		connstrs.PushBack(connstr.(string))
	}

	userDbLen := connstrs.Len()
	UserDb = make([]*Query, userDbLen)
	i := 0
	for e := connstrs.Front(); e != nil; e = e.Next() {
		UserDb[i] = &Query{}
		err := UserDb[i].Open(i+1, e.Value.(string))
		if err != nil {
			fmt.Printf("CANNOT OPEN USER DB %d: %s", i+1, err.Error())
		}
		i++
	}
}

func (p *Query) Open(dbId int, connstr string) error {
	var err error
	p.connString = connstr
	p.dbId = dbId
	p.db, err = sql.Open("mssql", connstr)
	if err != nil {
		fmt.Println("cannot open database:", err.Error())
		return err
	}

	return nil
}

func (p *Query) Close() {
	if p.db != nil {
		p.db.Close()
	}
}

func rowsToList(rows *sql.Rows) *list.List {
	results := list.New()

	cols, _ := rows.Columns()
	vals := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		vals[i] = new(interface{})
	}

	for rows.Next() {
		err := rows.Scan(vals...)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		row := make(map[string]interface{})

		for i := 0; i < len(vals); i++ {

			r := vals[i].(*interface{})

			switch v := (*r).(type) {
			case nil:
				row[cols[i]] = nil
			case bool:
				if v {
					row[cols[i]] = true
				} else {
					row[cols[i]] = false
				}
			case []byte:
				row[cols[i]] = string(v)
			default:
				row[cols[i]] = v
			}
		}
		results.PushBack(row)
	}

	return results
}

func (p *Query) Query(query string) *list.List {
	if p.db == nil {
		return nil
	}

	rows, err := p.db.Query(query)
	if err != nil {
		fmt.Println("cannot query:", err.Error())
		return nil
	}

	return rowsToList(rows)

}

func (p *Query) Proc(query string, args ...interface{}) *list.List {

	if len(args) > 0 {
		for i := 1; i <= len(args); i++ {
			if i > 1 {
				query += ","
			}
			query += " ?" + strconv.Itoa(i)
		}
	}

	stmt, err := p.db.Prepare(query)
	if err != nil {
		fmt.Println("Prepare:", err.Error())
		return nil
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		fmt.Println("Query:", err.Error())
		return nil
	}

	return rowsToList(rows)
}

func NewQuery() *Query {
	return &Query{}
}
