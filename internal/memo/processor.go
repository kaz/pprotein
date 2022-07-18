package memo

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/kaz/pprotein/internal/collect"
)

type (
	processor struct{}
)

func (p *processor) Cacheable() bool {
	return false
}

func (p *processor) Process(snapshot *collect.Snapshot) (io.ReadCloser, error) {
	bodyPath, err := snapshot.BodyPath()
	if err != nil {
		return nil, fmt.Errorf("failed to find snapshot body: %w", err)
	}

	res, err := ioutil.ReadFile(bodyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read snapshot body: %w", err)
	}
	return io.NopCloser(bytes.NewBuffer(res)), nil
}
