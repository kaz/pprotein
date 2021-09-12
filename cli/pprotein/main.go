package main

import (
	"net/http"
	"os"

	"github.com/kaz/pprotein/integration/echov4"
	"github.com/kaz/pprotein/internal/collect"
	"github.com/kaz/pprotein/internal/event"
	"github.com/kaz/pprotein/internal/extproc/alp"
	"github.com/kaz/pprotein/internal/extproc/querydigest"
	"github.com/kaz/pprotein/internal/pprof"
	"github.com/kaz/pprotein/view"
	"github.com/labstack/echo/v4"
)

func start() error {
	e := echo.New()
	echov4.Integrate(e)

	fs, err := view.FS()
	if err != nil {
		return err
	}
	e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(fs))))

	hub := event.NewHub()
	hub.RegisterHandlers(e.Group("/api/event"))

	pOpts := &collect.Options{
		WorkDir:   "./data/pprof",
		FileName:  "profile.pb.gz",
		EventName: "pprof",
		EventHub:  hub,
	}
	if err := pprof.RegisterHandlers(e.Group("/api/pprof"), pOpts); err != nil {
		return err
	}

	aOpts := &collect.Options{
		WorkDir:   "./data/httplog",
		FileName:  "raw.txt",
		EventName: "httplog",
		EventHub:  hub,
	}
	if err := alp.RegisterHandlers(e.Group("/api/httplog"), aOpts); err != nil {
		return err
	}

	qdOpts := &collect.Options{
		WorkDir:   "./data/slowlog",
		FileName:  "raw.txt",
		EventName: "slowlog",
		EventHub:  hub,
	}
	if err := querydigest.RegisterHandlers(e.Group("/api/slowlog"), qdOpts); err != nil {
		return err
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	return e.Start(":" + port)
}

func main() {
	if err := start(); err != nil {
		panic(err)
	}
}
