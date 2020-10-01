package main

import (
	"net/http"

	rice "github.com/GeertJohan/go.rice"
	"github.com/kaz/pprotein/embed"
	"github.com/kaz/pprotein/pprof"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	embed.EnableLogging(e)

	view := http.FileServer(rice.MustFindBox("view/dist").HTTPBox())
	e.GET("/*", echo.WrapHandler(view))

	pprofHandler, err := pprof.NewHandlers(pprof.Config{
		Workdir: "./tmp/pprof",
	})
	if err != nil {
		panic(err)
	}

	if err := pprofHandler.Register(e.Group("/api/pprof")); err != nil {
		panic(err)
	}

	panic(e.Start("0:9000"))
}
