package rest

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func NewServer() *echo.Echo {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		fmt.Println("hello there") // TODO remove
		return nil
	})

	e.GET("/latest", GetLatest())
	e.GET("/:date", GetByDate())

	return e
}

func GetLatest() func(c echo.Context) error {
	return func(c echo.Context) error {
		return nil
	}
}

func GetByDate() func(c echo.Context) error {
	return func(c echo.Context) error {
		date := c.Param("date")

		fmt.Println(date)

		return nil
	}
}
