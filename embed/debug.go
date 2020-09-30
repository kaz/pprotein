// +build !release

package main

import (
	"net/http"
	"net/http/pprof"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func echoPProf(e *echo.Echo) {
	pprofGroup := e.Group("/debug/pprof")
	pprofGroup.Any("/cmdline", echo.WrapHandler(http.HandlerFunc(pprof.Cmdline)))
	pprofGroup.Any("/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
	pprofGroup.Any("/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
	pprofGroup.Any("/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))
	pprofGroup.Any("/*", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
}

func echoLogging(e *echo.Echo) {
	e.Debug = true
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.Logger())
}
