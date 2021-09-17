package setting

import (
	"fmt"
	"io"
	"net/http"

	"github.com/kaz/pprotein/internal/storage"
	"github.com/labstack/echo/v4"
)

type (
	Handler struct {
		id       string
		defaults string

		store storage.Storage
	}
)

func NewHandler(id string, defaults string, store storage.Storage) *Handler {
	return &Handler{
		id:       id,
		defaults: defaults,
		store:    store,
	}
}

func (h *Handler) Register(g *echo.Group) {
	g.GET("", h.get)
	g.POST("", h.post)
}

func (h *Handler) get(c echo.Context) error {
	ok, err := h.store.HasBlob(h.id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to check blob state: %v", err))
	}
	if !ok {
		return c.String(http.StatusOK, h.defaults)
	}

	path, err := h.store.GetBlobPath(h.id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get blob: %v", err))
	}
	return c.File(path)
}
func (h *Handler) post(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to read body: %v", err))
	}

	if err := h.store.PutBlob(h.id, body); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to save: %v", err))
	}
	return c.NoContent(http.StatusNoContent)
}
