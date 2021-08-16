package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/kaz/kataribe"
	"github.com/kaz/pprotein/integration/echov4"
	"github.com/kaz/pprotein/internal/httplog"
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

	pprofCfg := pprof.Config{
		Workdir: "./data/pprof",
	}
	if err := pprof.RegisterHandlers(e.Group("/api/pprof"), pprofCfg); err != nil {
		return err
	}

	kataribeCfg := kataribe.Config{}
	if _, err := toml.DecodeFile("./kataribe.toml", &kataribeCfg); err != nil {
		return err
	}

	httplogCfg := httplog.Config{
		Workdir:  "./data/httplog",
		Kataribe: kataribeCfg,
	}
	if err := httplog.RegisterHandlers(e.Group("/api/httplog"), httplogCfg); err != nil {
		return err
	}

	slowlogCfg := slowlog.Config{
		Workdir: "./data/slowlog",
	}
	if err := slowlog.RegisterHandlers(e.Group("/api/slowlog"), slowlogCfg); err != nil {
		return err
	}

	return e.Start(":9000")
}

func main() {
	if err := start(); err != nil {
		panic(err)
	}
}
