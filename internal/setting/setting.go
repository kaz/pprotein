package setting

import (
	"fmt"
	"io"
	"net/http"

	"github.com/kaz/pprotein/internal/storage"
	"github.com/labstack/echo/v4"
)

type (
	handler struct {
		store storage.Storage

		id   string
		Path string
	}
)

func NewHandler(store storage.Storage, id string, defaults []byte) (*handler, error) {
	ok, err := store.ExistsFile(id)
	if err != nil {
		return nil, fmt.Errorf("failed to check blob state: %v", err)
	}
	if !ok {
		if err := store.PutAsFile(id, defaults); err != nil {
			return nil, fmt.Errorf("failed to save default value: %v", err)
		}
	}

	path, err := store.GetFilePath(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find blob path: %v", err)
	}

	return &handler{
		store: store,
		id:    id,
		Path:  path,
	}, nil
}

func (h *handler) Register(g *echo.Group) {
	g.GET("", h.get)
	g.POST("", h.post)
}

func (h *handler) get(c echo.Context) error {
	return c.File(h.Path)
}
func (h *handler) post(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to read body: %v", err))
	}

	pretty, err := sanitize(h.id, body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to parse body: %v", err))
	}

	if err := h.store.PutAsFile(h.id, pretty); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to save: %v", err))
	}

	return c.NoContent(http.StatusNoContent)
}
