package persistent

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/kaz/pprotein/internal/storage"
	"github.com/labstack/echo/v4"
)

type (
	Handler struct {
		store storage.Storage

		sanitize func([]byte) ([]byte, error)

		fileName string
		filePath string
	}
)

func New(store storage.Storage, fileName string, defaultContent []byte, sanitize func([]byte) ([]byte, error)) (*Handler, error) {
	ok, err := store.ExistsFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to check blob state: %v", err)
	}
	if !ok {
		content, err := sanitize(defaultContent)
		if err != nil {
			return nil, fmt.Errorf("failed to sanitize: %v", err)
		}
		if err := store.PutFile(fileName, content); err != nil {
			return nil, fmt.Errorf("failed to save default value: %v", err)
		}
	}

	filePath, err := store.GetFilePath(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to find blob path: %v", err)
	}

	return &Handler{
		store: store,

		sanitize: sanitize,

		fileName: fileName,
		filePath: filePath,
	}, nil
}

func (h *Handler) RegisterHandlers(g *echo.Group) {
	g.GET("", h.handleGet)
	g.POST("", h.handlePost)
}

func (h *Handler) GetPath() string {
	return h.filePath
}
func (h *Handler) GetContent() ([]byte, error) {
	content, err := os.ReadFile(h.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return content, nil
}

func (h *Handler) handleGet(c echo.Context) error {
	return c.File(h.filePath)
}
func (h *Handler) handlePost(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to read body: %v", err))
	}

	pretty, err := h.sanitize(body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to parse body: %v", err))
	}

	if err := h.store.PutFile(h.fileName, pretty); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to save: %v", err))
	}

	return c.NoContent(http.StatusOK)
}
