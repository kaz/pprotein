package kataribe

import (
	"fmt"
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
)

func RegisterHandlers(g *echo.Group, config Config) error {
	confPath, err := generateConfig(config)
	if err != nil {
		return fmt.Errorf("failed to generate kataribe config: %w", err)
	}

	p := extproc.NewProcessor("kataribe", "-conf", confPath)

	if err := extproc.RegisterHandlers(g, extproc.Config{Workdir: config.Workdir, Processor: p}); err != nil {
		return fmt.Errorf("failed to register extproc handlers: %w", err)
	}

	return nil
}

func generateConfig(config Config) (string, error) {
	confPath := path.Join(config.Workdir, "kataribe.toml")

	if _, err := os.Stat(confPath); err == nil {
		return confPath, nil
	}

	if err := os.MkdirAll(config.Workdir, 0755); err != nil {
		return "", fmt.Errorf("failed to make directory: %w", err)
	}
	if err := exec.Command("kataribe", "-conf", confPath, "-generate").Run(); err != nil {
		return "", fmt.Errorf("failed to run external process: %w", err)
	}
	return confPath, nil
}
