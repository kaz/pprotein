package httplog

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/kaz/kataribe"
	"github.com/kaz/pprotein/fetch"
	"github.com/labstack/echo"
)

type (
	Config struct {
		Workdir  string
		Kataribe kataribe.Config
	}
	Handler struct {
		store *fetch.Store
		proc  *processor
	}
	processor struct {
		store  sync.Map
		config kataribe.Config
	}
)

func RegisterHandlers(g *echo.Group, config Config) error {
	manager, err := fetch.NewManager(config.Workdir, "log")
	if err != nil {
		return fmt.Errorf("failed to initialize manager: %w", err)
	}

	proc := &processor{config: config.Kataribe}
	store, err := fetch.NewStore(manager, proc.process)
	if err != nil {
		return fmt.Errorf("failed to initialize store: %w", err)
	}

	h := &Handler{store: store, proc: proc}
	g.GET("/logs", h.logsGet)
	g.POST("/logs", h.logsPost)
	g.GET("/logs/:id", h.logsIdGet)

	return store.RegisterHandlers(g)
}

func (p *processor) process(e *fetch.Entry) error {
	file, err := os.Open(e.Path())
	if err != nil {
		return fmt.Errorf("failed to open log: %w", err)
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if err := kataribe.New(file, p.config).Print(buf); err != nil {
		return fmt.Errorf("failed to parse log: %w", err)
	}

	p.store.Store(e.ID, buf.Bytes())
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
	return c.Blob(http.StatusOK, "text/plain", data.([]byte))
}
