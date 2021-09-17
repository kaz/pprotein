package querydigest

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"

	"github.com/kaz/pprotein/internal/collect"
)

type (
	processor struct{}
)

func (p *processor) Cacheable() bool {
	return true
}

func (p *processor) Process(snapshot *collect.Snapshot) (io.ReadCloser, error) {
	bodyPath, err := snapshot.BodyPath()
	if err != nil {
		return nil, fmt.Errorf("failed to find snapshot body: %w", err)
	}

	cmd := exec.Command("pt-query-digest", "--limit", "100%", "--output", "json", bodyPath)

	res, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("external process aborted: %w", err)
	}

	starts := bytes.IndexByte(res, '{')
	return io.NopCloser(bytes.NewBuffer(res[starts:])), nil
}
