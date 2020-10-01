// +build !release

package embed

import (
	"net/http"
	"net/http/pprof"

	"github.com/kaz/pprotein/embed/transport"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func EnableLogTransport(e *echo.Echo) {
	g := e.Group("/debug/log")
	g.Any("/app", transport.NewTailHandler("./app.log").Handle)
	g.Any("/nginx", transport.NewTailHandler("/var/log/nginx/access.log").Handle)
}

func EnablePProf(e *echo.Echo) {
	g := e.Group("/debug/pprof")
	g.Any("/cmdline", echo.WrapHandler(http.HandlerFunc(pprof.Cmdline)))
	g.Any("/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
	g.Any("/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
	g.Any("/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))
	g.Any("/*", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
}

func EnableLogging(e *echo.Echo) {
	e.Debug = true
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.Logger())
}
