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
	"github.com/kaz/pprotein/internal/setting"
	"github.com/kaz/pprotein/internal/storage"
	"github.com/kaz/pprotein/view"
	"github.com/labstack/echo/v4"
)

func start() error {
	e := echo.New()
	echov4.Integrate(e)

	store, err := storage.New("data")
	if err != nil {
		return err
	}

	fs, err := view.FS()
	if err != nil {
		return err
	}
	e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(fs))))

	setting.NewHandler("group.json", "TODO", store).Register(e.Group("/api/setting/group"))
	setting.NewHandler("repository.json", "TODO", store).Register(e.Group("/api/setting/repository"))
	setting.NewHandler("httplog_config.yml", "TODO", store).Register(e.Group("/api/setting/httplog"))

	hub := event.NewHub()
	hub.RegisterHandlers(e.Group("/api/event"))

	pOpts := &collect.Options{
		Type:     "pprof",
		Ext:      "-pprof.pb.gz",
		Store:    store,
		EventHub: hub,
	}
	if err := pprof.RegisterHandlers(e.Group("/api/pprof"), pOpts); err != nil {
		return err
	}

	aOpts := &collect.Options{
		Type:     "httplog",
		Ext:      "-httplog.log",
		Store:    store,
		EventHub: hub,
	}
	if err := alp.RegisterHandlers(e.Group("/api/httplog"), aOpts); err != nil {
		return err
	}

	qdOpts := &collect.Options{
		Type:     "slowlog",
		Ext:      "-slowlog.log",
		Store:    store,
		EventHub: hub,
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
