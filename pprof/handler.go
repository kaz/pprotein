package pprof

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/google/pprof/driver"
	"github.com/kaz/pprotein/fetch"
	"github.com/labstack/echo/v4"
)

type (
	Config struct {
		Workdir string
	}
	Handler struct {
		store *fetch.Store
	}
	processor struct {
		route *echo.Group
		mu    sync.Mutex
	}
)

func RegisterHandlers(g *echo.Group, config Config) error {
	manager, err := fetch.NewManager(config.Workdir, "pb.gz")
	if err != nil {
		return fmt.Errorf("failed to initialize manager: %w", err)
	}

	proc := &processor{route: g}
	store, err := fetch.NewStore(manager, proc.process)
	if err != nil {
		return fmt.Errorf("failed to initialize store: %w", err)
	}

	h := &Handler{store: store}
	g.GET("", h.profilesGet)
	g.POST("", h.profilesPost)

	return store.RegisterHandlers(g)
}

func (p *processor) process(e *fetch.Entry) error {
	registerProfileHandlers := func(args *driver.HTTPServerArgs) error {
		if args.Hostport != "0:0" {
			return fmt.Errorf("unxpected hostport: %v", args.Hostport)
		}

		p.mu.Lock()
		defer p.mu.Unlock()

		ig := p.route.Group(getProfilePath(e))
		for key, handler := range args.Handlers {
			ig.Any(key, echo.WrapHandler(handler))
		}
		return nil
	}

	options := &driver.Options{
		Flagset: NewFlagSet([]string{
			"-no_browser",
			"-http", "0:0",
			e.Path(),
		}),
		HTTPServer: registerProfileHandlers,
	}

	if err := driver.PProf(options); err != nil {
		return fmt.Errorf("pprof internal error: %w", err)
	}
	return nil
}
func getProfilePath(e *fetch.Entry) string {
	return fmt.Sprintf("/profiles/%s", e.ID)
}

func (h *Handler) profilesGet(c echo.Context) error {
	return c.JSON(http.StatusOK, h.store.Get())
}

func (h *Handler) profilesPost(c echo.Context) error {
	req := &fetch.AddRequest{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to parse request body: %v", err))
	}

	if err := h.store.Add(req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to fetch profile: %v", err))
	}

	return c.NoContent(http.StatusAccepted)
}
