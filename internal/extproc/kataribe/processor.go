package kataribe

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/kaz/pprotein/internal/collect"
	"github.com/labstack/gommon/log"
)

type (
	processor struct {
		confPath string
	}
)

func (p *processor) Cacheable() bool {
	return true
}

func (p *processor) Process(snapshot *collect.Snapshot) (io.ReadCloser, error) {
	body, err := os.Open(snapshot.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	cmd := exec.Command("kataribe", "-conf", p.confPath)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to open stdout: %w", err)
	}

	go func() {
		defer stdin.Close()

		if _, err := io.Copy(stdin, body); err != nil {
			log.Error("[!] failed to write stdin of external process:", err.Error())
		}
	}()

	res, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("external process aborted: %w", err)
	}

	return io.NopCloser(bytes.NewBuffer(res)), nil
}
