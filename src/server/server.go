package main

import (
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
	title := c.FormValue("title")
	author := c.FormValue("author")

	dbbooks.InsertData(title, author, table)
	return c.JSON(http.StatusOK, &ResultResponse{Message: "Create done"})
}

func read(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, dbbooks.SelectData(id, table))
}

func update(c echo.Context) error {
	bookMap := echo.Map{}
	if e := c.Bind(&bookMap); e != nil {
		panic(e.Error())
	}

	id, _ := strconv.Atoi(bookMap["ID"].(string))
	dbbooks.UpdateData(id, bookMap["Title"].(string), bookMap["Author"].(string), table)
	// bookMap["Message"] = "Update done"
	return c.JSON(http.StatusOK, &ResultResponse{Message: "Update done"})
}

func delete(c echo.Context) error {
	auth := echo.Map{}
	if e := c.Bind(&auth); e != nil {
		panic(e.Error())
	}

	id, _ := strconv.Atoi(c.Param("id"))
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

	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/books", index)

	e.POST("/books", create)
	e.GET("/books/:id", read)
	e.PUT("/books", update)
	e.DELETE("/books/:id", delete)

	e.Logger.Fatal(e.Start(":1323"))
}
