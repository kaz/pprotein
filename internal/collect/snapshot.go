package collect

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

type (
	Snapshot struct {
		SnapshotMeta
		SnapshotPath `json:"-"`
	}
	SnapshotMeta struct {
		ID       string
		Datetime time.Time
		URL      string
		Duration int
	}
	SnapshotPath struct {
		Meta  string
		Body  string
		Cache string
	}
)

func (s *Snapshot) Prune() error {
	if err := os.RemoveAll(path.Dir(s.Body)); err != nil {
		return fmt.Errorf("failed to remove directory: %w", err)
	}
	return nil
}

func (s *Snapshot) Collect() error {
	resp, err := http.Get(fmt.Sprintf("%s?seconds=%d", s.URL, s.Duration))
	if err != nil {
		return fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	if err := os.Mkdir(path.Dir(s.Body), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	bodyFile, err := os.Create(s.Body)
	if err != nil {
		return fmt.Errorf("failed to create body file: %w", err)
	}
	defer bodyFile.Close()

	metaFile, err := os.Create(s.Meta)
	if err != nil {
		return fmt.Errorf("failed to create meta file: %w", err)
	}
	defer metaFile.Close()

	if _, err := io.Copy(bodyFile, resp.Body); err != nil {
		return fmt.Errorf("failed to write body: %w", err)
	}

	if err := json.NewEncoder(metaFile).Encode(s.SnapshotMeta); err != nil {
		return fmt.Errorf("failed to write meta: %w", err)
	}

	return nil
}
