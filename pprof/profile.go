package pprof

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type (
	Profiler struct {
		workdir string
	}

	Profile struct {
		bodyPath string

		ID       string
		Datetime time.Time
		URL      string
		Duration int
	}
)

func NewProfiler(workdir string) (*Profiler, error) {
	if err := os.MkdirAll(workdir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create workdir: %w", err)
	}
	return &Profiler{workdir: workdir}, nil
}

func (pr *Profiler) Fetch(url string, duration int, result chan func() (*Profile, error)) {
	res, err := pr.fetch(url, duration)
	result <- func() (*Profile, error) { return res, err }
}
func (pr *Profiler) fetch(url string, duration int) (*Profile, error) {
	ts := time.Now()
	id := fmt.Sprintf("%d", ts.UnixNano())

	bodyPath := path.Join(pr.workdir, fmt.Sprintf("%s.pb.gz", id))
	metaPath := path.Join(pr.workdir, fmt.Sprintf("%s.json", id))

	prof := &Profile{
		bodyPath: bodyPath,
		ID:       id,
		Datetime: ts,
		URL:      url,
		Duration: duration,
	}

	bodyFile, err := os.Create(bodyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create body file: %w", err)
	}
	defer bodyFile.Close()

	metaFile, err := os.Create(metaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create meta file: %w", err)
	}
	defer metaFile.Close()

	resp, err := http.Get(fmt.Sprintf("%s?second=%d", url, duration))
	if err != nil {
		return nil, fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	if _, err := io.Copy(bodyFile, resp.Body); err != nil {
		return nil, fmt.Errorf("failed to write body: %w", err)
	}

	if err := json.NewEncoder(metaFile).Encode(prof); err != nil {
		return nil, fmt.Errorf("failed to write meta: %w", err)
	}

	return prof, nil
}

func (pr *Profiler) List() ([]*Profile, error) {
	entries, err := ioutil.ReadDir(pr.workdir)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %w", err)
	}

	profiles := []*Profile{}
	for _, ent := range entries {
		if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".json") {
			continue
		}

		bodyPath := path.Join(pr.workdir, strings.Replace(ent.Name(), ".json", ".pb.gz", 1))
		metaPath := path.Join(pr.workdir, ent.Name())

		metaFile, err := os.Open(metaPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open meta file: %w", err)
		}
		defer metaFile.Close()

		prof := &Profile{bodyPath: bodyPath}
		if err := json.NewDecoder(metaFile).Decode(prof); err != nil {
			return nil, fmt.Errorf("failed to decode meta: %w", err)
		}

		profiles = append(profiles, prof)
	}

	return profiles, nil
}

func (p *Profile) Path() string {
	return p.bodyPath
}
