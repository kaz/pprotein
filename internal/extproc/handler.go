package extproc

import (
	"fmt"
	"net/http"

	"github.com/kaz/pprotein/internal/collect"
	"github.com/labstack/echo/v4"
)

type (
	Config struct {
		Workdir   string
		Processor collect.Processor
	}
	handler struct {
		collector *collect.Collector
	}
)

func RegisterHandlers(g *echo.Group, config Config) error {
	collector, err := collect.New(config.Processor, config.Workdir, "log.txt")
	if err != nil {
		return fmt.Errorf("failed to initialize collector: %w", err)
	}

	h := &handler{collector: collector}
	g.GET("", h.getIndex)
	g.POST("", h.postIndex)
	g.GET("/:id", h.getId)

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
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to fetch log: %v", err))
	}

	return c.NoContent(http.StatusAccepted)
}

func (h *handler) getId(c echo.Context) error {
	r, err := h.collector.Get(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get entry: %w", err))
	}
	defer r.Close()

	return c.Stream(http.StatusOK, "application/json", r)
}
