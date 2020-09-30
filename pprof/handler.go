package pprof

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/google/pprof/driver"
	"github.com/labstack/echo"
)

const (
	ProfileStatusOk      ProfileStatus = "ok"
	ProfileStatusFail    ProfileStatus = "fail"
	ProfileStatusPending ProfileStatus = "pending"
)

type (
	Config struct {
		Workdir string
	}
	Handler struct {
		profiler *Profiler
		stats    sync.Map

		route   *echo.Group
		routeMu sync.Mutex
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
	h.route = g

	g.GET("/profiles", h.profilesGet)
	g.POST("/profiles", h.profilesPost)

	profiles, err := h.profiler.List()
	if err != nil {
		return fmt.Errorf("failed to read profiles: %w", err)
	}

	for _, p := range profiles {
		h.parseProfile(p)
	}

	return nil
}

func (h *Handler) parseProfile(p *Profile) {
	h.stats.Store(p.ID, &ProfileInfo{
		Status:  ProfileStatusPending,
		Message: "Parsing",
		Profile: p,
	})

	ch := make(chan error)
	go func() {
		ch <- h.parseProfileSync(p)
	}()
	go func() {
		if err := <-ch; err != nil {
			h.stats.Store(p.ID, &ProfileInfo{
				Status:  ProfileStatusFail,
				Message: err.Error(),
				Profile: p,
			})
			return
		}
		h.stats.Store(p.ID, &ProfileInfo{
			Status:  ProfileStatusOk,
			Message: "",
			Profile: p,
		})
	}()
}
func (h *Handler) parseProfileSync(p *Profile) error {
	registerProfileHandlers := func(args *driver.HTTPServerArgs) error {
		if args.Hostport != "0:0" {
			return fmt.Errorf("unxpected hostport: %v", args.Hostport)
		}

		h.routeMu.Lock()
		defer h.routeMu.Unlock()

		ig := h.route.Group(getProfilePath(p))
		for key, handler := range args.Handlers {
			ig.Any(key, echo.WrapHandler(handler))
		}
		return nil
	}

	options := &driver.Options{
		Flagset: NewFlagSet([]string{
			"-no_browser",
			"-http", "0:0",
			p.Path(),
		}),
		HTTPServer: registerProfileHandlers,
	}

	if err := driver.PProf(options); err != nil {
		return fmt.Errorf("pprof internal error: %w", err)
	}
	return nil
}
func getProfilePath(p *Profile) string {
	return fmt.Sprintf("/profiles/%s", p.ID)
}

func (h *Handler) profilesGet(c echo.Context) error {
	data := map[string]interface{}{}
	h.stats.Range(func(key, value interface{}) bool {
		data[key.(string)] = value
		return true
	})
	return c.JSON(http.StatusOK, data)
}

func (h *Handler) profilesPost(c echo.Context) error {
	req := &FetchRequest{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to parse request body: %v", err))
	}

	if req.URL == "" || req.Duration == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "unexpected request format")
	}

	p := h.profiler.CreateProfile(req.URL, req.Duration)
	h.stats.Store(p.ID, &ProfileInfo{
		Status:  ProfileStatusPending,
		Message: "Fetching",
		Profile: p,
	})

	ch := make(chan error)
	go p.Fetch(ch)
	go func() {
		if err := <-ch; err != nil {
			h.stats.Store(p.ID, &ProfileInfo{
				Status:  ProfileStatusFail,
				Message: err.Error(),
				Profile: p,
			})
			return
		}
		h.parseProfile(p)
	}()

	return c.NoContent(http.StatusNoContent)
}
