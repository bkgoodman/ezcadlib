package main

import ( 
	"testing"
	"fmt"
	 //"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/jmoiron/sqlx"
)

func TestDB(t *testing.T) {
	db, _ := sqlx.Open("sqlite3", "./params.db")
	rows, _ := db.Queryx("SELECT * FROM params;")
	fmt.Printf("ROWS %+v\n",rows)
	cols, ee := rows.Columns()
	fmt.Printf("COLUMNS %s %q\n",ee,cols)
	for rows.Next() {
		fmt.Printf("ROW %+v\n",rows)
		var p params
		err:=rows.StructScan(&p)
		fmt.Printf("SCANROW %s %+v\n",err,p)
	}
}
