package querydigest

import (
	"fmt"

	"github.com/kaz/pprotein/internal/extproc"
	"github.com/labstack/echo/v4"
)

type (
	Config struct {
		Workdir string
	}
)

func RegisterHandlers(g *echo.Group, config Config) error {
	p := &processor{}
	if err := extproc.RegisterHandlers(g, extproc.Config{Workdir: config.Workdir, Processor: p}); err != nil {
		return fmt.Errorf("failed to register extproc handlers: %w", err)
	}

	return nil
}
