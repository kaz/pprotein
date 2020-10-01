package transport

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
	seconds, err := strconv.Atoi(c.QueryParam("seconds"))
	if err != nil {
		seconds = 30
	}

	t, err := tail.TailFile(h.filename, tail.Config{
		Follow:    true,
		ReOpen:    true,
		MustExist: true,
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: os.SEEK_END,
		},
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to start tailing: %v", err))
	}

	r, w := io.Pipe()
	timer := time.NewTimer(time.Duration(seconds) * time.Second)

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
