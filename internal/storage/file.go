package storage

import (
	"fmt"
	"os"
	"path"
)

type (
	fileStore struct {
		workdir string
	}
)

func newFile(workdir string) (fileStorage, error) {
	if err := os.MkdirAll(workdir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create workdir: %w", err)
	}

	return &fileStore{workdir}, nil
}

func (s *fileStore) PutFile(id string, data []byte) error {
	file, err := os.Create(path.Join(s.workdir, id))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
func (s *fileStore) GetFilePath(id string) (string, error) {
	return path.Join(s.workdir, id), nil
}
func (s *fileStore) ExistsFile(id string) (bool, error) {
	_, err := os.Stat(path.Join(s.workdir, id))
	return err == nil, nil
}
