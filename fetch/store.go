package fetch

import (
	"fmt"
	"sync"
)

const (
	StatusOk      Status = "ok"
	StatusFail    Status = "fail"
	StatusPending Status = "pending"
)

type (
	Store struct {
		manager   *Manager
		processor ProcessFn

		store sync.Map
	}

	ProcessFn func(*Entry) error

	Status string

	EntryInfo struct {
		Status  Status
		Message string
		Entry   *Entry
	}

	AddRequest struct {
		URL      string
		Duration int
	}
)

func NewStore(manager *Manager, processor ProcessFn) (*Store, error) {
	s := &Store{
		manager:   manager,
		processor: processor,
	}

	entries, err := s.manager.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list entries: %w", err)
	}

	for _, e := range entries {
		s.parse(e)
	}

	return s, nil
}

func (s *Store) updateStatus(e *Entry, status Status, msg string) {
	s.store.Store(e.ID, &EntryInfo{
		Status:  status,
		Message: msg,
		Entry:   e,
	})
}

func (s *Store) parse(e *Entry) {
	s.updateStatus(e, StatusPending, "Parsing")

	ch := make(chan error)
	go func() {
		ch <- s.processor(e)
	}()
	go func() {
		if err := <-ch; err != nil {
			s.updateStatus(e, StatusFail, err.Error())
			return
		}
		s.updateStatus(e, StatusOk, "Ready")
	}()
}

func (s *Store) Get() map[string]*EntryInfo {
	data := map[string]*EntryInfo{}
	s.store.Range(func(key, value interface{}) bool {
		data[key.(string)] = value.(*EntryInfo)
		return true
	})
	return data
}

func (s *Store) Add(req *AddRequest) error {
	if req.URL == "" || req.Duration == 0 {
		return fmt.Errorf("any parameters cannot be nil")
	}

	e := s.manager.Create(req.URL, req.Duration)
	s.updateStatus(e, StatusPending, "Collecting")

	ch := make(chan error)
	go e.Fetch(ch)
	go func() {
		if err := <-ch; err != nil {
			s.updateStatus(e, StatusFail, err.Error())
			return
		}
		s.parse(e)
	}()

	return nil
}
