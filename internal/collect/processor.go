package collect

import (
	"fmt"
	"io"
	"os"
)

type (
	Processor interface {
		Process(snapshot *Snapshot) (io.ReadCloser, error)
		Cacheable() bool
	}

	cachedProcessor struct {
		internal Processor
	}
)

func newCachedProcessor(internal Processor) Processor {
	return &cachedProcessor{internal}
}

func (p *cachedProcessor) Process(snapshot *Snapshot) (io.ReadCloser, error) {
	if _, err := os.Stat(snapshot.Cache); err == nil {
		return p.serveCached(snapshot)
	}
	return p.serveGenerated(snapshot)
}
func (p *cachedProcessor) serveCached(snapshot *Snapshot) (io.ReadCloser, error) {
	cache, err := os.Open(snapshot.Cache)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	return cache, nil
}
func (p *cachedProcessor) serveGenerated(snapshot *Snapshot) (io.ReadCloser, error) {
	r, err := p.internal.Process(snapshot)
	if err != nil {
		return nil, fmt.Errorf("internal processor returns an error: %w", err)
	}

	if !p.internal.Cacheable() {
		return r, nil
	}

	cache, err := os.Create(snapshot.Cache)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	if _, err := io.Copy(cache, r); err != nil {
		return nil, fmt.Errorf("failed to copy: %w", err)
	}
	cache.Close()

	return p.serveCached(snapshot)
}

func (p *cachedProcessor) Cacheable() bool {
	return false
}
