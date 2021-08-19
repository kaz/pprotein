package kataribe

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/kaz/pprotein/internal/extproc"
	"github.com/labstack/echo/v4"
)

type (
	Config struct {
		Workdir string
	}
	Handler struct {
		confPath string
	}
)

func RegisterHandlers(g *echo.Group, config Config) error {
	h, err := newHandler(config)
	if err != nil {
		return fmt.Errorf("failed to initialize handler: %w", err)
	}

	p := &processor{confPath: h.confPath}
	if err := extproc.RegisterHandlers(g, extproc.Config{Workdir: config.Workdir, Processor: p}); err != nil {
		return fmt.Errorf("failed to register extproc handlers: %w", err)
	}

	g.GET("/config", h.getConfig)
	g.POST("/config", h.postConfig)

	return nil
}

func newHandler(config Config) (*Handler, error) {
	h := &Handler{confPath: path.Join(config.Workdir, "kataribe.toml")}
	if _, err := os.Stat(h.confPath); err == nil {
		return h, nil
	}

	if err := os.MkdirAll(config.Workdir, 0755); err != nil {
		return nil, fmt.Errorf("failed to make directory: %w", err)
	}
	if err := exec.Command("kataribe", "-conf", h.confPath, "-generate").Run(); err != nil {
		return nil, fmt.Errorf("failed to generate kataribe config: %w", err)
	}
	return h, nil
}

func (h *Handler) getConfig(c echo.Context) error {
	return c.File(h.confPath)
}
func (h *Handler) postConfig(c echo.Context) error {
	file, err := os.Create(h.confPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to create file: %v", err))
	}
	defer file.Close()

	if _, err := io.Copy(file, c.Request().Body); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to write to file: %v", err))
	}

	return c.NoContent(http.StatusNoContent)
}
