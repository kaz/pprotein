package storage

import "fmt"

type (
	store struct {
		kvStorage
		fileStorage
	}
)

func New(workdir string) (Storage, error) {
	kvs, err := newKV(workdir)
	if err != nil {
		return nil, fmt.Errorf("failed to create kvs: %w", err)
	}

	fs, err := newFile(workdir)
	if err != nil {
		return nil, fmt.Errorf("failed to create fs: %w", err)
	}

	return &store{kvs, fs}, nil
}
