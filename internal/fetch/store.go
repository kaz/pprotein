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

		mu    *sync.RWMutex
		data  map[string]*EntryInfo
		event *sync.Cond
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

		mu:   &sync.RWMutex{},
		data: map[string]*EntryInfo{},

		event: sync.NewCond(&sync.Mutex{}),
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
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[e.ID] = &EntryInfo{
		Status:  status,
		Message: msg,
		Entry:   e,
	}

	s.event.Broadcast()
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
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.data
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
