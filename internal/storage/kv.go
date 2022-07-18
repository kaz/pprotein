package storage

import (
	"fmt"
	"os"
	"path"

	"go.etcd.io/bbolt"
)

type (
	kvStore struct {
		db *bbolt.DB
	}
)

func newKV(workdir string) (kvStorage, error) {
	if err := os.MkdirAll(workdir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create workdir: %w", err)
	}

	db, err := bbolt.Open(path.Join(workdir, "pprotein.db"), 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize DB: %w", err)
	}
	return &kvStore{db}, nil
}

func (s *kvStore) Put(typ, id string, data []byte) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(typ))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		return bucket.Put([]byte(id), data)
	})
}
func (s *kvStore) Get(typ, id string) (resp []byte, err error) {
	err = s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(typ))
		if bucket == nil {
			return nil
		}

		resp = bucket.Get([]byte(id))
		return nil
	})
	return resp, err
}
func (s *kvStore) GetAll(typ string) (resp [][]byte, err error) {
	err = s.db.View(func(tx *bbolt.Tx) error {
		resp = make([][]byte, 0)

		bucket := tx.Bucket([]byte(typ))
		if bucket == nil {
			return nil
		}

		bucket.ForEach(func(k, v []byte) error {
			resp = append(resp, v)
			return nil
		})
		return nil
	})
	return resp, err
}
func (s *kvStore) Exists(typ, id string) (exists bool, err error) {
	err = s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(typ))
		if bucket == nil {
			return nil
		}

		exists = bucket.Get([]byte(id)) != nil
		return nil
	})
	return exists, err
}
func (s *kvStore) Delete(typ, id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(typ))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		return bucket.Delete([]byte(id))
	})
}
