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
		collector *collect.Collector
	}
)

func RegisterHandlers(g *echo.Group, opts *collect.Options) error {
	p := &processor{mu: &sync.Mutex{}, route: g}

	collector, err := collect.New(p, opts)
	if err != nil {
		return fmt.Errorf("failed to initialize collector: %w", err)
	}

	h := &handler{collector: collector}
	g.GET("", h.getIndex)
	g.POST("", h.postIndex)

	return nil
}

func (h *handler) getIndex(c echo.Context) error {
	return c.JSON(http.StatusOK, h.collector.List())
}

func (h *handler) postIndex(c echo.Context) error {
	req := &collect.Job{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to parse request body: %v", err))
	}

	go func() {
		if err := h.collector.Collect(req); err != nil {
			log.Error("[!] collector aborted:", err)
		}
	}()

	return c.NoContent(http.StatusAccepted)
}
