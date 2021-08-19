package collect

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"
)

type (
	Storage struct {
		workdir  string
		filename string
	}
)

func newStorage(workdir string, filename string) (*Storage, error) {
	if err := os.MkdirAll(workdir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create workdir: %w", err)
	}

	return &Storage{
		workdir:  workdir,
		filename: filename,
	}, nil
}

func (s *Storage) getPathById(id string) SnapshotPath {
	return SnapshotPath{
		Meta:  path.Join(s.workdir, id, "index.json"),
		Body:  path.Join(s.workdir, id, s.filename),
		Cache: path.Join(s.workdir, id, "cache.dat"),
	}
}

func (s *Storage) PrepareSnapshot(url string, duration int) *Snapshot {
	ts := time.Now()
	id := strconv.FormatInt(ts.UnixNano(), 36)

	return &Snapshot{
		SnapshotMeta: SnapshotMeta{
			ID:       id,
			Datetime: ts,
			URL:      url,
			Duration: duration,
		},
		SnapshotPath: s.getPathById(id),
	}
}

func (s *Storage) List() ([]*Snapshot, error) {
	finfos, err := ioutil.ReadDir(s.workdir)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %w", err)
	}

	snapshots := []*Snapshot{}
	for _, finfo := range finfos {
		if !finfo.IsDir() {
			continue
		}

		id := path.Base(finfo.Name())
		sPath := s.getPathById(id)

		metaFile, err := os.Open(sPath.Meta)
		if err != nil {
			go (&Snapshot{SnapshotPath: sPath}).Prune()
			fmt.Fprintf(os.Stderr, "[!] ignored=%v: failed to open meta file: %v\n", id, err)
			continue
		}
		defer metaFile.Close()

		sMeta := SnapshotMeta{}
		if err := json.NewDecoder(metaFile).Decode(&sMeta); err != nil {
			go (&Snapshot{SnapshotPath: sPath}).Prune()
			fmt.Fprintf(os.Stderr, "[!] ignored=%v: failed to decode meta file: %v\n", id, err)
			continue
		}

		snapshots = append(snapshots, &Snapshot{SnapshotMeta: sMeta, SnapshotPath: sPath})
	}

	return snapshots, nil
}
