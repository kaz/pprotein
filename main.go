package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/kaz/pprotein/integration/echov4"
	"github.com/kaz/pprotein/internal/extproc/kataribe"
	"github.com/kaz/pprotein/internal/pprof"
	"github.com/kaz/pprotein/internal/slowlog"
	"github.com/labstack/echo/v4"
)

//go:embed view/dist/*
var view embed.FS

func start() error {
	e := echo.New()
	echov4.Integrate(e)

	subfs, err := fs.Sub(view, "view/dist")
	if err != nil {
		return err
	}
	e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(subfs))))

	if err := pprof.RegisterHandlers(e.Group("/api/pprof"), pprof.Config{Workdir: "./data/pprof"}); err != nil {
		return err
	}
	if err := kataribe.RegisterHandlers(e.Group("/api/httplog"), kataribe.Config{Workdir: "./data/httplog"}); err != nil {
		return err
	}
	if err := slowlog.RegisterHandlers(e.Group("/api/slowlog"), slowlog.Config{Workdir: "./data/slowlog"}); err != nil {
		return err
	}

	return e.Start(":9000")
}

func main() {
	if err := start(); err != nil {
		panic(err)
	}
}
