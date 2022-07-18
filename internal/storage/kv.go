package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path"

	_ "modernc.org/sqlite"
)

type (
	kvStore struct {
		db *sql.DB
	}
)

var (
	initsqls = []string{
		"CREATE TABLE IF NOT EXISTS snapshot (type TEXT, id TEXT, content BLOB, deleted INTEGER, PRIMARY KEY(type, id))",
		"CREATE INDEX IF NOT EXISTS snapshot_type_deleted ON snapshot (type, deleted)",
		"CREATE INDEX IF NOT EXISTS snapshot_type_id_deleted ON snapshot (type, id, deleted)",
	}
)

func newKV(workdir string) (kvStorage, error) {
	if err := os.MkdirAll(workdir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create workdir: %w", err)
	}

	db, err := sql.Open("sqlite", path.Join(workdir, "kv.sqlite3"))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize DB: %w", err)
	}

	for _, sql := range initsqls {
		if _, err := db.Exec(sql); err != nil {
			return nil, fmt.Errorf("failed to run init sql: %w", err)
		}
	}
	return &kvStore{db}, nil
}

func (s *kvStore) Put(typ, id string, data []byte) error {
	if _, err := s.db.Exec("INSERT INTO snapshot VALUES (?, ?, ?, 0)", typ, id, data); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}
func (s *kvStore) Get(typ, id string) ([]byte, error) {
	data := []byte{}
	if err := s.db.QueryRow("SELECT content FROM snapshot WHERE type = ? AND id = ? AND deleted = 0", typ, id).Scan(&data); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return data, nil
}
func (s *kvStore) GetAll(typ string) ([][]byte, error) {
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
func (s *kvStore) Exists(typ, id string) (bool, error) {
	cnt := 0
	if err := s.db.QueryRow("SELECT COUNT(*) FROM snapshot WHERE type = ? AND id = ?", typ, id).Scan(&cnt); err != nil {
		return false, fmt.Errorf("failed to select: %w", err)
	}
	return cnt != 0, nil
}
func (s *kvStore) Delete(typ, id string) error {
	if _, err := s.db.Exec("UPDATE snapshot SET deleted = 1 WHERE type = ? AND id = ?", typ, id); err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}
	return nil
}
