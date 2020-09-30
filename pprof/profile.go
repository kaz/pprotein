package pprof

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
	Profiler struct {
		workdir string
	}

	Profile struct {
		bodyPath string
		metaPath string

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

func (pr *Profiler) CreateProfile(url string, duration int) *Profile {
	ts := time.Now()
	id := strconv.FormatInt(ts.UnixNano(), 36)

	return &Profile{
		bodyPath: path.Join(pr.workdir, fmt.Sprintf("%s.pb.gz", id)),
		metaPath: path.Join(pr.workdir, fmt.Sprintf("%s.json", id)),

		ID:       id,
		Datetime: ts,
		URL:      url,
		Duration: duration,
	}
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
			fmt.Fprintf(os.Stderr, "[SKIPPED] failed to open meta file: %v", err)
			continue
		}
		defer metaFile.Close()

		prof := &Profile{
			bodyPath: bodyPath,
			metaPath: metaPath,
		}
		if err := json.NewDecoder(metaFile).Decode(prof); err != nil {
			fmt.Fprintf(os.Stderr, "[SKIPPED] failed to decode meta: %v", err)
			continue
		}

		profiles = append(profiles, prof)
	}

	return profiles, nil
}

func (p *Profile) Path() string {
	return p.bodyPath
}

func (p *Profile) Fetch(result chan error) {
	result <- p.fetch()
}
func (p *Profile) fetch() error {
	resp, err := http.Get(fmt.Sprintf("%s?second=%d", p.URL, p.Duration))
	if err != nil {
		return fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	bodyFile, err := os.Create(p.bodyPath)
	if err != nil {
		return fmt.Errorf("failed to create body file: %w", err)
	}
	defer bodyFile.Close()

	metaFile, err := os.Create(p.metaPath)
	if err != nil {
		return fmt.Errorf("failed to create meta file: %w", err)
	}
	defer metaFile.Close()

	if _, err := io.Copy(bodyFile, resp.Body); err != nil {
		return fmt.Errorf("failed to write body: %w", err)
	}

	if err := json.NewEncoder(metaFile).Encode(p); err != nil {
		return fmt.Errorf("failed to write meta: %w", err)
	}

	return nil
}
