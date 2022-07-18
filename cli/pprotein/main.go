package main

import (
	_ "embed"
	"net/http"
	"os"

	"github.com/kaz/pprotein/integration/echov4"
	"github.com/kaz/pprotein/internal/collect"
	"github.com/kaz/pprotein/internal/event"
	"github.com/kaz/pprotein/internal/extproc/alp"
	"github.com/kaz/pprotein/internal/extproc/querydigest"
	"github.com/kaz/pprotein/internal/memo"
	"github.com/kaz/pprotein/internal/pprof"
	"github.com/kaz/pprotein/internal/setting"
	"github.com/kaz/pprotein/internal/storage"
	"github.com/kaz/pprotein/view"
	"github.com/labstack/echo/v4"
)

//go:embed group.json
var defaultGroupJson []byte

//go:embed alp.yml
var defaultAlpYml []byte

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

	groupJson, err := setting.NewHandler(store, "group.json", defaultGroupJson)
	if err != nil {
		return err
	}
	groupJson.Register(e.Group("/api/setting/group"))

	alpYml, err := setting.NewHandler(store, "alp.yml", defaultAlpYml)
	if err != nil {
		return err
	}
	alpYml.Register(e.Group("/api/setting/httplog"))

	hub := event.NewHub()
	hub.RegisterHandlers(e.Group("/api/event"))

	pOpts := &collect.Options{
		Type:     "pprof",
		Ext:      "-pprof.pb.gz",
		Store:    store,
		EventHub: hub,
	}
	if err := pprof.NewHandler(pOpts).Register(e.Group("/api/pprof")); err != nil {
		return err
	}

	aOpts := &collect.Options{
		Type:     "httplog",
		Ext:      "-httplog.log",
		Store:    store,
		EventHub: hub,
	}
	if err := alp.NewHandler(alpYml.Path, aOpts).Register(e.Group("/api/httplog")); err != nil {
		return err
	}

	qOpts := &collect.Options{
		Type:     "slowlog",
		Ext:      "-slowlog.log",
		Store:    store,
		EventHub: hub,
	}
	if err := querydigest.NewHandler(qOpts).Register(e.Group("/api/slowlog")); err != nil {
		return err
	}

	mOpts := &collect.Options{
		Type:     "memo",
		Ext:      "-memo.log",
		Store:    store,
		EventHub: hub,
	}
	if err := memo.NewHandler(mOpts).Register(e.Group("/api/memo")); err != nil {
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
