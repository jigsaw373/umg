package application

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func setMiddlewares(e *echo.Echo) {
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		ExposeHeaders: []string{"X-Pagination-Page-Count"},
	}))
}
