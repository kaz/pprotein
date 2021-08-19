package event

import (
	"github.com/alexandrevicenzi/go-sse"
	"github.com/labstack/echo/v4"
)

type (
	Hub struct {
		server *sse.Server
	}
)

func NewHub() *Hub {
	return &Hub{
		server: sse.NewServer(&sse.Options{}),
	}
}

func (h *Hub) RegisterHandlers(g *echo.Group) {
	g.GET("", echo.WrapHandler(h.server))
}

func (h *Hub) Publisher(message string) *Publisher {
	return &Publisher{
		hub:     h,
		message: message,
	}
}
func (h *Hub) publish(message string) {
	h.server.SendMessage("", sse.SimpleMessage(string(message)))
}
