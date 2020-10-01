package main

import (
	"net/http"

	"github.com/BurntSushi/toml"
	rice "github.com/GeertJohan/go.rice"
	"github.com/kaz/kataribe"
	"github.com/kaz/pprotein/embed"
	"github.com/kaz/pprotein/httplog"
	"github.com/kaz/pprotein/pprof"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	embed.EnableLogging(e)

	view := http.FileServer(rice.MustFindBox("view/dist").HTTPBox())
	e.GET("/*", echo.WrapHandler(view))

	pprofCfg := pprof.Config{
		Workdir: "./data/pprof",
	}
	if err := pprof.RegisterHandlers(e.Group("/api/pprof"), pprofCfg); err != nil {
		panic(err)
	}

	kataribeCfg := kataribe.Config{}
	if _, err := toml.DecodeFile("./kataribe.toml", &kataribeCfg); err != nil {
		panic(err)
	}

	httplogCfg := httplog.Config{
		Workdir:  "./data/httplog",
		Kataribe: kataribeCfg,
	}
	if err := httplog.RegisterHandlers(e.Group("/api/httplog"), httplogCfg); err != nil {
		panic(err)
	}

	panic(e.Start(":9000"))
}
