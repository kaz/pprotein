package fetch

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func (s *Store) RegisterHandlers(g *echo.Group) error {
	g.GET("/events", s.eventsGet)
	return nil
}

func (s *Store) eventsGet(c echo.Context) error {
	r := c.Response()
	r.Header().Set("Content-Type", "text/event-stream")
	r.Header().Set("Cache-Control", "no-cache")
	r.Header().Set("Connection", "keep-alive")
	r.WriteHeader(http.StatusOK)
	r.Flush()

	for {
		select {
		case <-c.Request().Context().Done():
			return nil
		default:
			s.event.L.Lock()
			s.event.Wait()
			s.event.L.Unlock()
			fmt.Fprintf(r, "data: \n\n")
			r.Flush()
		}
	}
}
