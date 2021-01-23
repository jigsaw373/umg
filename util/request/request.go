package request

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetPagination(c echo.Context) (count, page int64) {
	// pagination parameters
	var err error

	count, err = strconv.ParseInt(c.QueryParam("count"), 10, 64)
	if err != nil || count < 1 || count > 20 {
		count = 20
	}

	page, err = strconv.ParseInt(c.QueryParam("page"), 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	return
}
