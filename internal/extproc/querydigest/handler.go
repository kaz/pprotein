package querydigest

import (
	"fmt"

	"github.com/kaz/pprotein/internal/collect"
	"github.com/kaz/pprotein/internal/extproc"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(g *echo.Group, opts *collect.Options) error {
	p := &processor{}
	if err := extproc.RegisterHandlers(g, p, opts); err != nil {
		return fmt.Errorf("failed to register extproc handlers: %w", err)
	}

	return nil
}
