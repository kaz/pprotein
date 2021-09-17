package querydigest

import (
	"fmt"

	"github.com/kaz/pprotein/internal/collect"
	"github.com/kaz/pprotein/internal/extproc"
	"github.com/labstack/echo/v4"
)

type (
	handler struct {
		opts *collect.Options
	}
)

func NewHandler(opts *collect.Options) *handler {
	return &handler{opts: opts}
}

func (h *handler) Register(g *echo.Group) error {
	if err := extproc.NewHandler(&processor{}, h.opts).Register(g); err != nil {
		return fmt.Errorf("failed to register extproc handlers: %w", err)
	}
	return nil
}
