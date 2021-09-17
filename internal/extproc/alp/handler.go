package alp

import (
	_ "embed"
	"fmt"

	"github.com/kaz/pprotein/internal/collect"
	"github.com/kaz/pprotein/internal/extproc"
	"github.com/labstack/echo/v4"
)

type (
	handler struct {
		confPath string
		opts     *collect.Options
	}
)

func NewHandler(confPath string, opts *collect.Options) *handler {
	return &handler{
		confPath: confPath,
		opts:     opts,
	}
}

func (h *handler) Register(g *echo.Group) error {
	if err := extproc.NewHandler(&processor{confPath: h.confPath}, h.opts).Register(g); err != nil {
		return fmt.Errorf("failed to register extproc handlers: %w", err)
	}
	return nil
}
