package collect

import (
	"bytes"
	"fmt"
	"io"

	"github.com/kaz/pprotein/internal/storage"
)

const cacheTypeKey = "cache"

type (
	Processor interface {
		Process(snapshot *Snapshot) (io.ReadCloser, error)
		Cacheable() bool
	}

	cachedProcessor struct {
		internal Processor
		store    storage.Storage
	}
)

func newCachedProcessor(internal Processor, store storage.Storage) Processor {
	return &cachedProcessor{internal, store}
}

func (p *cachedProcessor) Process(snapshot *Snapshot) (io.ReadCloser, error) {
	if ok, err := p.store.Exists(cacheTypeKey, snapshot.ID); err != nil {
		return nil, fmt.Errorf("failed to check cache status: %w", err)
	} else if ok {
		return p.serveCached(snapshot)
	}
	return p.serveGenerated(snapshot)
}
func (p *cachedProcessor) serveCached(snapshot *Snapshot) (io.ReadCloser, error) {
	cache, err := p.store.Get(cacheTypeKey, snapshot.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache: %w", err)
	}
	return io.NopCloser(bytes.NewBuffer(cache)), nil
}
func (p *cachedProcessor) serveGenerated(snapshot *Snapshot) (io.ReadCloser, error) {
	r, err := p.internal.Process(snapshot)
	if err != nil {
		return nil, fmt.Errorf("internal error: %w", err)
	}

	if !p.internal.Cacheable() {
		return r, nil
	}

	cacheContent, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	if err := p.store.Put(cacheTypeKey, snapshot.ID, cacheContent); err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}
	return p.serveCached(snapshot)
}

func (p *cachedProcessor) Cacheable() bool {
	return false
}
