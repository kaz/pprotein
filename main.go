package main

import (
	"github.com/kaz/pprotein/embed"
	"github.com/kaz/pprotein/pprof"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	embed.EnableLogging(e)

	pprofHandler, err := pprof.NewHandlers(pprof.Config{
		Workdir: "./tmp/pprof",
	})
	if err != nil {
		panic(err)
	}

	if err := pprofHandler.Register(e.Group("/pprof")); err != nil {
		panic(err)
	}

	panic(e.Start("0:9000"))
}
