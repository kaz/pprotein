package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type (
	Manager struct {
		workdir   string
		extention string
	}

	Entry struct {
		bodyPath string
		metaPath string

		ID       string
		Datetime time.Time
		URL      string
		Duration int
	}
)

func NewManager(workdir string, extention string) (*Manager, error) {
	if err := os.MkdirAll(workdir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create workdir: %w", err)
	}
	m := &Manager{
		workdir:   workdir,
		extention: extention,
	}
	return m, nil
}

func (m *Manager) Create(url string, duration int) *Entry {
	ts := time.Now()
	id := strconv.FormatInt(ts.UnixNano(), 36)

	return &Entry{
		bodyPath: path.Join(m.workdir, fmt.Sprintf("%s.%s", id, m.extention)),
		metaPath: path.Join(m.workdir, fmt.Sprintf("%s.json", id)),

		ID:       id,
		Datetime: ts,
		URL:      url,
		Duration: duration,
	}
}

func (m *Manager) List() ([]*Entry, error) {
	finfos, err := ioutil.ReadDir(m.workdir)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %w", err)
	}

	entries := []*Entry{}
	for _, finfo := range finfos {
		if finfo.IsDir() || !strings.HasSuffix(finfo.Name(), ".json") {
			continue
		}

		bodyPath := path.Join(m.workdir, strings.Replace(finfo.Name(), ".json", fmt.Sprintf(".%s", m.extention), 1))
		metaPath := path.Join(m.workdir, finfo.Name())

		metaFile, err := os.Open(metaPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[SKIPPED] failed to open meta file: %v", err)
			continue
		}
		defer metaFile.Close()

		ent := &Entry{
			bodyPath: bodyPath,
			metaPath: metaPath,
		}
		if err := json.NewDecoder(metaFile).Decode(ent); err != nil {
			fmt.Fprintf(os.Stderr, "[SKIPPED] failed to decode meta: %v", err)
			continue
		}

		entries = append(entries, ent)
	}

	return entries, nil
}

func (e *Entry) Path() string {
	return e.bodyPath
}

func (e *Entry) Fetch(result chan error) {
	result <- e.fetch()
}
func (e *Entry) fetch() error {
	resp, err := http.Get(fmt.Sprintf("%s?second=%d", e.URL, e.Duration))
	if err != nil {
		return fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	bodyFile, err := os.Create(e.bodyPath)
	if err != nil {
		return fmt.Errorf("failed to create body file: %w", err)
	}
	defer bodyFile.Close()

	metaFile, err := os.Create(e.metaPath)
	if err != nil {
		return fmt.Errorf("failed to create meta file: %w", err)
	}
	defer metaFile.Close()

	if _, err := io.Copy(bodyFile, resp.Body); err != nil {
		return fmt.Errorf("failed to write body: %w", err)
	}

	if err := json.NewEncoder(metaFile).Encode(e); err != nil {
		return fmt.Errorf("failed to write meta: %w", err)
	}

	return nil
}
