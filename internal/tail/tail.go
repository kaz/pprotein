package tail

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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
		log.Printf("serve failed: %v", err)
	}
}
func (h *TailHandler) serve(w http.ResponseWriter, r *http.Request) error {
	seconds, err := strconv.Atoi(r.URL.Query().Get("seconds"))
	if err != nil {
		seconds = 30
	}

	var output io.Writer = w
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		ew, err := gzip.NewWriterLevel(w, gzip.DefaultCompression)
		if err != nil {
			return fmt.Errorf("failed to initialize gzip writer: %w", err)
		}
		defer ew.Close()

		output = ew
		w.Header().Set("Content-Encoding", "gzip")
	}

	if err := h.tail(output, time.Duration(seconds)*time.Second); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		output.Write([]byte(err.Error()))

		return fmt.Errorf("failed to tail: %w", err)
	}

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	return nil
}

func (h *TailHandler) tail(w io.Writer, duration time.Duration) error {
	file, err := os.Open(h.filename)
	if err != nil {
		return fmt.Errorf("failed to open: %w", err)
	}
	defer file.Close()

	startPos, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("failed to seek: %w", err)
	}

	time.Sleep(duration)

	finfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat: %w", err)
	}

	r := io.LimitReader(file, finfo.Size()-startPos)
	if _, err := io.Copy(w, r); err != nil {
		return fmt.Errorf("failed to copy: %w", err)
	}
	return nil
}
