package pprof

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/kaz/pprotein/internal/collect"
	"github.com/labstack/echo/v4"
)

type (
	Config struct {
		Workdir string
	}
	handler struct {
		collector *collect.Collector
	}
)

func RegisterHandlers(g *echo.Group, config Config) error {
	p := &processor{mu: &sync.Mutex{}, route: g}

	collector, err := collect.New(p, config.Workdir, "profile.pb.gz")
	if err != nil {
		return fmt.Errorf("failed to initialize collector: %w", err)
	}

	h := &handler{collector: collector}
	g.GET("", h.getIndex)
	g.POST("", h.postIndex)

	return collector.RegisterHandlers(g)
}

func (h *handler) getIndex(c echo.Context) error {
	return c.JSON(http.StatusOK, h.collector.List())
}

func (h *handler) postIndex(c echo.Context) error {
	req := &collect.CollectRequest{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to parse request body: %v", err))
	}

	if err := h.collector.Collect(req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to fetch profile: %v", err))
	}

	return c.NoContent(http.StatusAccepted)
}
