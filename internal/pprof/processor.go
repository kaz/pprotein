package pprof

import (
	"fmt"
	"io"
	"sync"

	"github.com/google/pprof/driver"
	"github.com/kaz/pprotein/internal/collect"
	"github.com/labstack/echo/v4"
)

type (
	processor struct {
		mu    *sync.Mutex
		route *echo.Group
	}
)

func (p *processor) Cacheable() bool {
	return false
}

func (p *processor) Process(snapshot *collect.Snapshot) (io.ReadCloser, error) {
	registerProfileHandlers := func(args *driver.HTTPServerArgs) error {
		if args.Hostport != "0:0" {
			return fmt.Errorf("unxpected hostport: %v", args.Hostport)
		}

		p.mu.Lock()
		defer p.mu.Unlock()

		ig := p.route.Group(fmt.Sprintf("/%s", snapshot.ID))
		for key, handler := range args.Handlers {
			ig.Any(key, echo.WrapHandler(handler))
		}
		return nil
	}

	bodyPath, err := snapshot.BodyPath()
	if err != nil {
		return nil, fmt.Errorf("failed to find snapshot body: %w", err)
	}

	options := &driver.Options{
		Flagset: NewFlagSet([]string{
			"-no_browser",
			"-http", "0:0",
			bodyPath,
		}),
		HTTPServer: registerProfileHandlers,
	}

	if err := driver.PProf(options); err != nil {
		return nil, fmt.Errorf("pprof internal error: %w", err)
	}
	return nil, nil
}
