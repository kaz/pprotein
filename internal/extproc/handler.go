package extproc

import (
	"fmt"
	"net/http"

	"github.com/kaz/pprotein/internal/collect"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type (
	handler struct {
		processor collect.Processor
		opts      *collect.Options
		collector *collect.Collector
	}
)

func NewHandler(processor collect.Processor, opts *collect.Options) *handler {
	return &handler{
		processor: processor,
		opts:      opts,
	}
}

func (h *handler) Register(g *echo.Group) error {
	var err error
	h.collector, err = collect.New(h.processor, h.opts)
	if err != nil {
		return fmt.Errorf("failed to initialize collector: %w", err)
	}

	g.GET("", h.getIndex)
	g.POST("", h.postIndex)
	g.GET("/:id", h.getId)

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
