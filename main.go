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

	pprofCfg := pprof.Config{
		Workdir: "./tmp/pprof",
	}
	if err := pprof.RegisterHandlers(e.Group("/api/pprof"), pprofCfg); err != nil {
		panic(err)
	}

	panic(e.Start("0:9000"))
}
