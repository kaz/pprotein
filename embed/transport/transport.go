package transport

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/hpcloud/tail"
	"github.com/labstack/echo"
)

type (
	TailHandler struct {
		filename string
	}
)

func NewTailHandler(filename string) *TailHandler {
	return &TailHandler{filename}
}

func (h *TailHandler) Handle(c echo.Context) error {
	second, err := strconv.Atoi(c.QueryParam("second"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid `second` parameter: %v", err))
	}

	t, err := tail.TailFile(h.filename, tail.Config{
		Follow:    true,
		ReOpen:    true,
		MustExist: true,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to start tailing: %v", err))
	}

	r, w := io.Pipe()
	timer := time.NewTimer(time.Duration(second) * time.Second)

	go func() {
		defer r.Close()
		defer w.Close()
		defer t.Cleanup()

		for {
			select {
			case line := <-t.Lines:
				fmt.Fprintln(w, line.Text)
			case <-timer.C:
				return
			}
		}
	}()

	return c.Stream(http.StatusOK, "text/plain", r)
}
