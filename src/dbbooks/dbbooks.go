package dbbooks

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/syndtr/goleveldb/leveldb"
	_ "github.com/tidwall/buntdb"
)

// Book 구조체
type Book struct {
	ID     int
	Title  string
	Author string
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

	sql, err := db.Prepare("delete from `" + table + "` where _id=?")
	if err != nil {
		panic(err.Error())
	}
	sql.Exec(id)
}

func main() {
	// go run이나 debug로는 실행 안 됨. go build나 go install로 실행파일을 만들어서 실행해야 됨.
	// 실행파일(#1) 명령(#2) 대상(#3...)
	argLen := len(os.Args) - 1
	if argLen < 1 {
		fmt.Println("Usage:")
		fmt.Println("appmain.exe insert book_name author_name")
		fmt.Println("appmain.exe select [id_number]")
		fmt.Println("appmain.exe update id_number book_name author_name")
		fmt.Println("appmain.exe delete id_number")
		os.Exit(1)
		// panic("Please run with arguments.")
	}

	command := os.Args[1:2]

	table := "books"

	switch command[0] {
	case "insert":
		if argLen < 3 {
			fmt.Println(command[0] + ": Need more params")
			os.Exit(1)
		}
		bookName := os.Args[2:3]
		authorName := os.Args[3:4]

		InsertData(bookName[0], authorName[0], table)
	case "select":
		var id int

		if argLen < 1 {
			panic("Need more params")
		} else if argLen == 2 {
			idStr := os.Args[2:3]
			id, _ = strconv.Atoi(idStr[0])

		} else {
			id = 0
		}

		bookDatas := SelectData(id, table) // cRud - Read=Select
		for _, bookData := range bookDatas {
			fmt.Println(bookData.ID, bookData.Title, bookData.Author)
		}
	case "update":
		var id int

		if argLen < 4 {
			panic("Need more params")
		} else {
			idStr := os.Args[2:3]
			id, _ = strconv.Atoi(idStr[0])

			bookName := os.Args[3:4]
			authorName := os.Args[4:5]

			UpdateData(id, bookName[0], authorName[0], table) // crUd - Update
		}
	case "delete":
		var id int

		if argLen < 1 {
			panic("Need more params")
		} else if argLen == 2 {
			idStr := os.Args[2:3]
			id, _ = strconv.Atoi(idStr[0])

		} else {
			id = 0
		}

		if id > 0 {
			DeleteData(id, table) // cruD - Delete : 없는 _id를 골랐는데 에러가 안뜬다?? 머지머지??
		} else {
			panic("Id value have to be larger than 1.")
		}
	default:
		panic("Check inputed parameters.")
	}
}
