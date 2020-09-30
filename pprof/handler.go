package pprof

import (
	"fmt"
	"net/http"

	"github.com/google/pprof/driver"
	"github.com/labstack/echo"
)

type (
	Config struct {
		Workdir string
	}
	Handler struct {
		profiler *Profiler
	}

	FetchRequest struct {
		URL      string
		Duration int
	}
)

func NewHandlers(config Config) (*Handler, error) {
	profiler, err := NewProfiler(config.Workdir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize profiler: %w", err)
	}
	return &Handler{profiler}, nil
}

func (h *Handler) Register(g *echo.Group) error {
	g.POST("/fetch", h.fetch)

	profiles, err := h.profiler.List()
	if err != nil {
		return fmt.Errorf("failed to read profiles: %w", err)
	}

	for _, p := range profiles {
		registerProfile(g.Group(fmt.Sprintf("/%s", p.ID)), p)
	}

	return nil
}

func registerProfile(g *echo.Group, p *Profile) {
	ch := make(chan error)
	go func() {
		ch <- registerProfileSync(g, p)
	}()
	go func() {
		if err := <-ch; err != nil {
			// TODO CHANGE STATE
		}
		// TODO CHANGE STATE
	}()
}
func registerProfileSync(g *echo.Group, p *Profile) error {
	options := &driver.Options{
		Flagset: NewFlagSet([]string{
			"-no_browser",
			"-http", "0:0",
			p.Path(),
		}),
		HTTPServer: func(args *driver.HTTPServerArgs) error {
			if args.Hostport != "0:0" {
				return fmt.Errorf("unxpected hostport: %v", args.Hostport)
			}
			ig := g.Group(fmt.Sprintf("/%s", p.ID))
			for key, handler := range args.Handlers {
				ig.Any(key, echo.WrapHandler(handler))
			}
			return nil
		},
	}
	if err := driver.PProf(options); err != nil {
		return fmt.Errorf("pprof internal error: %w", err)
	}
	return nil
}

func (h *Handler) fetch(c echo.Context) error {
	g, ok := c.Get("group").(*echo.Group)
	if !ok {
		return fmt.Errorf("cannot retrieve router!")
	}

	req := &FetchRequest{}
	if err := c.Bind(req); err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}

	chRes := make(chan *Profile)
	chErr := make(chan error)

	go h.profiler.Fetch(req.URL, req.Duration, chRes, chErr)
	go func() {
		select {
		case err := <-chErr:
			if err != nil {
				// TODO CHANGE STATE
			}
		case p := <-chRes:
			// TODO CHANGE STATE
			registerProfile(g.Group(fmt.Sprintf("/%s", p.ID)), p)
		}
	}()

	return c.NoContent(http.StatusNoContent)
}
