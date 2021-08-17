package collect

import (
	"fmt"
	"io"
	"sync"
)

type (
	Collector struct {
		storage   *Storage
		processor Processor

		mu    *sync.RWMutex
		data  map[string]*Entry
		event *sync.Cond
	}

	Entry struct {
		Snapshot *Snapshot
		Status   Status
		Message  string
	}

	CollectRequest struct {
		URL      string
		Duration int
	}

	Status string
)

const (
	StatusOk      Status = "ok"
	StatusFail    Status = "fail"
	StatusPending Status = "pending"
)

func New(processor Processor, workdir string, filename string) (*Collector, error) {
	store, err := newStorage(workdir, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	c := &Collector{
		storage:   store,
		processor: newCachedProcessor(processor),

		mu:    &sync.RWMutex{},
		data:  map[string]*Entry{},
		event: sync.NewCond(&sync.Mutex{}),
	}

	snapshots, err := store.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list snapshots: %w", err)
	}

	for _, snapshot := range snapshots {
		go c.runProcessor(snapshot)
	}

	return c, nil
}

func (c *Collector) updateStatus(snapshot *Snapshot, status Status, msg string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[snapshot.ID] = &Entry{
		Snapshot: snapshot,
		Status:   status,
		Message:  msg,
	}

	c.event.Broadcast()
}

func (c *Collector) runProcessor(snapshot *Snapshot) error {
	c.updateStatus(snapshot, StatusPending, "Processing")

	r, err := c.processor.Process(snapshot)
	if err != nil {
		go snapshot.Prune()
		c.updateStatus(snapshot, StatusFail, err.Error())
		return fmt.Errorf("processor aborted: %w", err)
	}
	if r != nil {
		r.Close()
	}

	c.updateStatus(snapshot, StatusOk, "Ready")
	return nil
}

func (c *Collector) Get(id string) (io.ReadCloser, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ent, ok := c.data[id]
	if !ok {
		return nil, fmt.Errorf("no such entry: %v", ent)
	}

	return c.processor.Process(ent.Snapshot)
}

func (c *Collector) List() []*Entry {
	c.mu.RLock()
	defer c.mu.RUnlock()

	resp := make([]*Entry, 0, len(c.data))
	for _, ent := range c.data {
		resp = append(resp, ent)
	}
	return resp
}

func (c *Collector) Collect(req *CollectRequest) error {
	if req.URL == "" || req.Duration == 0 {
		return fmt.Errorf("any parameters cannot be nil")
	}

	snapshot := c.storage.PrepareSnapshot(req.URL, req.Duration)
	c.updateStatus(snapshot, StatusPending, "Collecting")

	if err := snapshot.Collect(); err != nil {
		c.updateStatus(snapshot, StatusFail, err.Error())
		return fmt.Errorf("failed to collect: %w", err)
	}

	if err := c.runProcessor(snapshot); err != nil {
		c.updateStatus(snapshot, StatusFail, err.Error())
		return fmt.Errorf("failed to process: %w", err)
	}
	return nil
}
