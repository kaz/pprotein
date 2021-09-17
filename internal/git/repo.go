package git

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/v39/github"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type (
	handler struct {
		confPath string
	}

	Config struct {
		Owner string
		Repo  string
		Token string
	}
)

func NewHandler(confPath string) *handler {
	return &handler{confPath: confPath}
}

func (h *handler) Register(g *echo.Group) {
	g.GET("/commit/:hash", h.getCommit)
}

func (h *handler) getCommit(c echo.Context) error {
	cfg, err := h.readConfig()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to read config: %v", err))
	}

	client := github.NewClient(oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.Token})))
	commit, _, err := client.Git.GetCommit(context.Background(), cfg.Owner, cfg.Repo, c.Param("hash"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get commit: %v", err))
	}

	return c.JSON(http.StatusOK, commit)
}

func (h *handler) readConfig() (*Config, error) {
	f, err := os.Open(h.confPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}
	return &cfg, nil
}
