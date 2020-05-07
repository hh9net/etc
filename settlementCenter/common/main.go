package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-oci8"
)

func main() {
	db, err := sql.Open("oci8", "admin/123@//192.168.0.160:1521/ORCL")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	rows, err := db.Query("select user from dual")
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		var data string
		rows.Scan(&data)
		fmt.Println(data)
	}

	if err := sqlQuery(db, "SELECT * FROM (SELECT t.*  FROM ADMIN.B_TXF_CHEDXFYSSJ t    ) WHERE ROWNUM <= 5"); err != nil {
		log.Fatal(err)
	}
	if err = rows.Err(); err != nil {
		log.Fatalln(err)
	}
}

//sql查询
func sqlQuery(db *sql.DB, sqlStmt string) error {
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return err
	}
	defer rows.Close()

	var n int
	for rows.Next() {
		fmt.Println("查询ok")
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	log.Printf("SQL Query success rows queried %d\n", n)
	return nil
}
