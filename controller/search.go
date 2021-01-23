package controller

import (
	"github.com/labstack/echo/v4"

	"github.com/boof/umg/services"
	"github.com/boof/umg/util/response"
)

func SearchUsers(c echo.Context) error {
	text := c.QueryParam("text")

	res := services.SearchUsers(text)
	return response.OK(c, res)
}

func SearchRoles(c echo.Context) error {
	text := c.QueryParam("text")

	res := services.SearchRoles(text)
	return response.OK(c, res)
}
