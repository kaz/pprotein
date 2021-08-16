package tail

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/hpcloud/tail"
)

type (
	TailHandler struct {
		filename string
	}
)

func NewTailHandler(filename string) *TailHandler {
	return &TailHandler{filename}
}

func (h *TailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.serve(w, r); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}
func (h *TailHandler) serve(w http.ResponseWriter, r *http.Request) error {
	seconds, err := strconv.Atoi(r.URL.Query().Get("seconds"))
	if err != nil {
		seconds = 30
	}

	f, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("type conversion failed")
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
		return fmt.Errorf("failed to start tail: %w", err)
	}
	defer t.Cleanup()

	timer := time.NewTimer(time.Duration(seconds) * time.Second)
	for {
		select {
		case <-timer.C:
			return nil
		case line := <-t.Lines:
			fmt.Fprintln(w, line.Text)
			f.Flush()
		}
	}
}
