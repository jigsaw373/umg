package application

import (
	"os"

	"github.com/labstack/echo/v4"
)

func RunServer() {
	e := echo.New()

	setMiddlewares(e)
	mapRoutes(e)

	e.Logger.Fatal(e.Start(os.Getenv("API_PORT")))
}
