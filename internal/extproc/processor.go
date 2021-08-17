package extproc

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/kaz/pprotein/internal/collect"
)

type (
	processor struct {
		name string
		args []string
	}
)

func NewProcessor(name string, args ...string) *processor {
	return &processor{name: name, args: args}
}

func (p *processor) Cacheable() bool {
	return true
}

func (p *processor) Process(snapshot *collect.Snapshot) (io.ReadCloser, error) {
	args := make([]string, 0, len(p.args)+1)
	args = append(args, p.args...)
	args = append(args, snapshot.Body)

	cmd := exec.Command(p.name, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to open stdout: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start external process: %w", err)
	}

	return stdout, nil
}
