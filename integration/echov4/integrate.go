package echov4

import (
	"github.com/kaz/pprotein/integration"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func Integrate(e *echo.Echo) {
	EnableDebugHandler(e)
	EnableDebugMode(e)
}

func EnableDebugHandler(e *echo.Echo) {
	e.Any("/debug/*", echo.WrapHandler(integration.NewDebugHandler()))
}

func EnableDebugMode(e *echo.Echo) {
	e.Debug = true
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
}
