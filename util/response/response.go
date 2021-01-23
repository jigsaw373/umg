package response

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func SetPageCountHeader(c *echo.Context, pages int64) {
	// set pagination header
	(*c).Response().Header().Add("X-Pagination-Page-Count", strconv.FormatInt(pages, 10))
}

// BadReq sends an error message with HTTP BAD REQUEST status code
func BadReq(c echo.Context, msg string) error {
	return c.JSON(http.StatusBadRequest, echo.Map{"message": msg})
}

// InternalErr sends and error message with HTTP INTERNAL ERROR status code
func InternalErr(c echo.Context, msg string) error {
	return c.JSON(http.StatusInternalServerError, echo.Map{"message": msg})
}

// NotFound sends not found message with HTTP NOT FOUND status code
func NotFound(c echo.Context, msg string) error {
	return c.JSON(http.StatusNotFound, echo.Map{"message": msg})
}

// Unauthorized send unauthorized message with HTTP UNAUTHORIZED status code
func Unauthorized(c echo.Context, msg string) error {
	return c.JSON(http.StatusUnauthorized, echo.Map{"message": msg})
}

func NotAcceptable(c echo.Context, msg string) error {
	return c.JSON(http.StatusNotAcceptable, echo.Map{"message": msg})
}

func Done(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{})
}

func OK(c echo.Context, result interface{}) error {
	return c.JSON(http.StatusOK, echo.Map{"result": result})
}

func Created(c echo.Context, result interface{}) error {
	return c.JSON(http.StatusCreated, echo.Map{"result": result})
}

func NotContent(c echo.Context) error {
	return c.JSON(http.StatusNoContent, echo.Map{})
}
