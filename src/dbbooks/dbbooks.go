package dbbooks

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// Book 구조체
type Book struct {
	ID     int    `json:"ID"`
	Title  string `json:"Title"`
	Author string `json:"Author"`
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "myslimsite"
	dbHost := "tcp(localhost:13306)"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@"+dbHost+"/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

// InsertData : Crud
func InsertData(title string, author string, table string) {
	db := dbConn()
	defer db.Close()

	sql, err := db.Prepare("insert into " + table + "(title, author) values(?, ?)")
	if err != nil {
		panic(err.Error())
	}
	sql.Exec(title, author)
}

// SelectData : cRud
func SelectData(id int, table string) (result []Book) {
	db := dbConn()
	defer db.Close()

	sql := "select * from `" + table + "`"
	var where string
	if id > 0 {
		where = " where `_id`='" + strconv.Itoa(id) + "'"
	} else {
		where = ""
	}
	order := " order by '_id' desc"

	sql = sql + where + order

	selDB, err := db.Query(sql)
	if err != nil {
		panic(err.Error())
	}

	book := Book{}
	result = []Book{}
	for selDB.Next() {
		var id int
		var title, author string
		err = selDB.Scan(&id, &title, &author)
		if err != nil {
			panic(err.Error())
		}
		book.ID = id
		book.Title = title
		book.Author = author
		result = append(result, book)
	}

	return
}

// UpdateData : crUd
func UpdateData(id int, title string, author string, table string) {
	db := dbConn()
	defer db.Close()

	sql, err := db.Prepare("update `" + table + "` set title=?, author=? where _id=?")
	if err != nil {
		panic(err.Error())
	}
	sql.Exec(title, author, id)
}

// DeleteData : cruD
func DeleteData(id int, table string) {
	db := dbConn()
	defer db.Close()

	fmt.Println(id, table)

	sql, err := db.Prepare("delete from `" + table + "` where `_id`=?")
	if err != nil {
		panic(err.Error())
	}
	sql.Exec(id)
}
