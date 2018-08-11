package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"dbbooks"
)

var table = "books"

// ResultResponse : Create, Read 결과 반환용
type ResultResponse struct{ Message string }

func index(c echo.Context) error {
	return c.JSON(http.StatusOK, dbbooks.SelectData(0, table))
}

func create(c echo.Context) error {
	title := c.FormValue("Title")
	author := c.FormValue("Author")

	dbbooks.InsertData(title, author, table)
	return c.JSON(http.StatusOK, &ResultResponse{Message: "Create done"})
}

func read(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, dbbooks.SelectData(id, table))
}

func update(c echo.Context) error {
	book := new(dbbooks.Book)
	if e := c.Bind(book); e != nil {
		panic(e.Error())
	}

	id, _ := strconv.Atoi(c.Param("id"))
	dbbooks.UpdateData(id, book.Title, book.Author, table)

	return c.JSON(http.StatusOK, &ResultResponse{Message: "Update done"})
}

func delete(c echo.Context) error {
	// Error is occured so, auth is waived.
	// auth := echo.Map{}
	// if e := c.Bind(&auth); e != nil {
	// 	panic(e.Error())
	// }

	id, _ := strconv.Atoi(c.Param("id"))

	fmt.Println(id)

	dbbooks.DeleteData(id, table)
	return c.JSON(http.StatusOK, &ResultResponse{Message: "Delete done"})
}

func main() {
	echo.NotFoundHandler = func(c echo.Context) error {
		errorResult := &ResultResponse{
			Message: "Contents not found",
		}
		return c.JSON(http.StatusNotFound, errorResult)
	}

	e := echo.New()

	e.Use(middleware.CORS())
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"http://localhost:8080", "*"},
	// 	AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	// }))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/books", index)

	e.POST("/books", create)
	e.GET("/books/:id", read)
	e.PUT("/books/:id", update)
	e.DELETE("/books/:id", delete)

	e.Logger.Fatal(e.Start(":1323"))
}
