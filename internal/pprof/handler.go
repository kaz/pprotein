package pprof

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/kaz/pprotein/internal/collect"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type (
	handler struct {
		opts      *collect.Options
		collector *collect.Collector
	}
)

func NewHandler(opts *collect.Options) *handler {
	return &handler{opts: opts}
}

func (h *handler) Register(g *echo.Group) error {
	p := &processor{mu: &sync.Mutex{}, route: g}

	var err error
	h.collector, err = collect.New(p, h.opts)
	if err != nil {
		return fmt.Errorf("failed to initialize collector: %w", err)
	}

	g.GET("", h.getIndex)
	g.POST("", h.postIndex)

	return nil
}

func (h *handler) getIndex(c echo.Context) error {
	return c.JSON(http.StatusOK, h.collector.List())
}

func (h *handler) postIndex(c echo.Context) error {
	target := &collect.SnapshotTarget{}
	if err := c.Bind(target); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to parse request body: %v", err))
	}

	go func() {
		if err := h.collector.Collect(target); err != nil {
			log.Error("[!] collector aborted:", err)
		}
	}()

	return c.NoContent(http.StatusOK)
}
