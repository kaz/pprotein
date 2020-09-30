package pprof

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/google/pprof/driver"
	"github.com/labstack/echo"
)

const (
	INJECTED_ROUTER_KEY = "router"

	STATUS_OK      ProfileStatus = "ok"
	STATUS_FAIL    ProfileStatus = "fail"
	STATUS_PENDING ProfileStatus = "pending"
)

type (
	Config struct {
		Workdir string
	}
	Handler struct {
		profiler *Profiler
		stats    sync.Map
	}

	ProfileStatus string
	ProfileInfo   struct {
		Status  ProfileStatus
		Message string
		Profile *Profile
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
	return &Handler{profiler: profiler}, nil
}

func (h *Handler) Register(g *echo.Group) error {
	g.Use(injectRouter(g))

	g.GET("/stats", h.statsHandler)
	g.POST("/fetch", h.fetchHandler)

	profiles, err := h.profiler.List()
	if err != nil {
		return fmt.Errorf("failed to read profiles: %w", err)
	}

	for _, p := range profiles {
		h.parseProfile(g.Group(fmt.Sprintf("/%s", p.ID)), p)
	}

	return nil
}

func injectRouter(g *echo.Group) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(INJECTED_ROUTER_KEY, g)
			return next(c)
		}
	}
}

func (h *Handler) parseProfile(g *echo.Group, p *Profile) {
	h.stats.Store(p.ID, &ProfileInfo{
		Status:  STATUS_PENDING,
		Message: "Parsing",
		Profile: p,
	})

	ch := make(chan error)
	go func() {
		ch <- registerProfileHandlers(g, p)
	}()
	go func() {
		if err := <-ch; err != nil {
			h.stats.Store(p.ID, &ProfileInfo{
				Status:  STATUS_FAIL,
				Message: err.Error(),
				Profile: p,
			})
			return
		}
		h.stats.Store(p.ID, &ProfileInfo{
			Status:  STATUS_OK,
			Message: "",
			Profile: p,
		})
	}()
}
func registerProfileHandlers(g *echo.Group, p *Profile) error {
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

func (h *Handler) statsHandler(c echo.Context) error {
	data := map[string]interface{}{}
	h.stats.Range(func(key, value interface{}) bool {
		data[key.(string)] = value
		return true
	})
	return c.JSON(http.StatusOK, data)
}

func (h *Handler) fetchHandler(c echo.Context) error {
	g, ok := c.Get(INJECTED_ROUTER_KEY).(*echo.Group)
	if !ok {
		return fmt.Errorf("cannot retrieve router!")
	}

	req := &FetchRequest{}
	if err := c.Bind(req); err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}

	p := h.profiler.CreateProfile(req.URL, req.Duration)
	h.stats.Store(p.ID, &ProfileInfo{
		Status:  STATUS_PENDING,
		Message: "Fetching",
		Profile: p,
	})

	ch := make(chan error)
	go p.Fetch(ch)
	go func() {
		if err := <-ch; err != nil {
			h.stats.Store(p.ID, &ProfileInfo{
				Status:  STATUS_FAIL,
				Message: err.Error(),
				Profile: p,
			})
			return
		}
		h.parseProfile(g.Group(fmt.Sprintf("/%s", p.ID)), p)
	}()

	return c.NoContent(http.StatusNoContent)
}
