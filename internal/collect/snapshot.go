package collect

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type (
	Snapshot struct {
		*SnapshotMeta
		*SnapshotTarget
		*SnapshotPath `json:"-"`
	}
	SnapshotMeta struct {
		Type        string
		ID          string
		Datetime    time.Time
		GitRevision string
	}
	SnapshotTarget struct {
		GroupId  string
		Label    string
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
	if err := os.Rename(s.Meta, s.Meta+".ignored"); err != nil {
		return fmt.Errorf("failed to rename meta file: %w", err)
	}
	return nil
}

func (s *Snapshot) Collect() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?seconds=%d", s.URL, s.Duration), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	var r io.Reader = resp.Body
	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		cr, err := gzip.NewReader(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to initialize gzip reader: %w", err)
		}
		defer cr.Close()

		r = cr
	}

	s.GitRevision = resp.Header.Get("X-GIT-REVISION")

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(r)
		return fmt.Errorf("http error: status=%v, body=%v", resp.StatusCode, string(body))
	}

	if err := os.Mkdir(path.Dir(s.Meta), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	metaFile, err := os.Create(s.Meta)
	if err != nil {
		return fmt.Errorf("failed to create meta file: %w", err)
	}
	defer metaFile.Close()

	bodyFile, err := os.Create(s.Body)
	if err != nil {
		return fmt.Errorf("failed to create body file: %w", err)
	}
	defer bodyFile.Close()

	if written, err := io.Copy(bodyFile, r); err != nil {
		return fmt.Errorf("failed to write body: %w", err)
	} else if written == 0 {
		return fmt.Errorf("empty response")
	}

	if err := json.NewEncoder(metaFile).Encode(s); err != nil {
		return fmt.Errorf("failed to write meta: %w", err)
	}

	return nil
}
