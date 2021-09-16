package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path"

	_ "modernc.org/sqlite"
)

type (
	storageImpl struct {
		db      *sql.DB
		workdir string
	}
)

var (
	initsqls = []string{
		"CREATE TABLE IF NOT EXISTS snapshot (id TEXT PRIMARY KEY, type TEXT, content BLOB, deleted INTEGER)",
		"CREATE TABLE IF NOT EXISTS cache (id TEXT PRIMARY KEY, content BLOB)",
		"CREATE INDEX IF NOT EXISTS snapshot_lookup ON snapshot (type, deleted)",
	}
)

func New(workdir string) (Storage, error) {
	if err := os.MkdirAll(workdir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create workdir: %w", err)
	}

	db, err := sql.Open("sqlite", path.Join(workdir, "data.sqlite3"))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize DB: %w", err)
	}

	for _, sql := range initsqls {
		if _, err := db.Exec(sql); err != nil {
			return nil, fmt.Errorf("failed to run init sql: %w", err)
		}
	}
	return &storageImpl{db, workdir}, nil
}

func (s *storageImpl) PutSnapshot(id string, typ string, data []byte) error {
	if _, err := s.db.Exec("INSERT INTO snapshot VALUES (?, ?, ?, 0)", id, typ, data); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}
func (s *storageImpl) PutCache(id string, data []byte) error {
	if _, err := s.db.Exec("INSERT INTO cache VALUES (?, ?)", id, data); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}
func (s *storageImpl) PutBlob(id string, data []byte) error {
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

func (s *storageImpl) GetSnapshots(typ string) ([][]byte, error) {
	rows, err := s.db.Query("SELECT content FROM snapshot WHERE type = ? AND deleted = 0", typ)
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	result := [][]byte{}
	for rows.Next() {
		data := []byte{}
		if err := rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		result = append(result, data)
	}

	return result, nil
}
func (s *storageImpl) GetCacheContent(id string) ([]byte, error) {
	data := []byte{}
	if err := s.db.QueryRow("SELECT content FROM cache WHERE id = ?", id).Scan(&data); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return data, nil
}
func (s *storageImpl) GetBlobPath(id string) (string, error) {
	return path.Join(s.workdir, id), nil
}

func (s *storageImpl) DeleteSnapshot(id string) error {
	if _, err := s.db.Exec("UPDATE snapshot SET deleted = 1 WHERE id = ?", id); err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}
	return nil
}
func (s *storageImpl) HasCache(id string) (bool, error) {
	cnt := 0
	if err := s.db.QueryRow("SELECT COUNT(1) FROM cache WHERE id = ?", id).Scan(&cnt); err != nil {
		return false, fmt.Errorf("failed to select: %w", err)
	}
	return cnt != 0, nil
}
