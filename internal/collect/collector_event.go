package collect

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (col *Collector) RegisterHandlers(g *echo.Group) error {
	g.GET("/events", col.eventsGet)
	return nil
}

func (col *Collector) eventsGet(c echo.Context) error {
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
			col.event.L.Lock()
			col.event.Wait()
			col.event.L.Unlock()
			fmt.Fprintf(r, "data: \n\n")
			r.Flush()
		}
	}
}
