package slowlog

import (
	"fmt"
	"net/http"
	"os/exec"
	"sync"

	"github.com/kaz/pprotein/internal/fetch"
	"github.com/labstack/echo/v4"
)

type (
	Config struct {
		Workdir string
	}
	Handler struct {
		store *fetch.Store
		proc  *processor
	}
	processor struct {
		store sync.Map
	}
)

func RegisterHandlers(g *echo.Group, config Config) error {
	manager, err := fetch.NewManager(config.Workdir, "log")
	if err != nil {
		return fmt.Errorf("failed to initialize manager: %w", err)
	}

	proc := &processor{}
	store, err := fetch.NewStore(manager, proc.process)
	if err != nil {
		return fmt.Errorf("failed to initialize store: %w", err)
	}

	h := &Handler{store: store, proc: proc}
	g.GET("", h.logsGet)
	g.POST("", h.logsPost)
	g.GET("/:id", h.logsIdGet)

	return store.RegisterHandlers(g)
}

func (p *processor) process(e *fetch.Entry) error {
	cmd := exec.Command("pt-query-digest", "--limit", "100%", "--output", "json", e.Path())
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to exec pt-query-digest: %w", err)
	}

	p.store.Store(e.ID, out)
	return nil
}

func (h *Handler) logsGet(c echo.Context) error {
	return c.JSON(http.StatusOK, h.store.Get())
}

func (h *Handler) logsPost(c echo.Context) error {
	req := &fetch.AddRequest{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to parse request body: %v", err))
	}

	if err := h.store.Add(req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to fetch log: %v", err))
	}

	return c.NoContent(http.StatusAccepted)
}

func (h *Handler) logsIdGet(c echo.Context) error {
	data, ok := h.proc.store.Load(c.Param("id"))
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return c.Blob(http.StatusOK, "application/json", data.([]byte))
}
