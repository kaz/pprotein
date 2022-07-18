package setting

import (
	"fmt"
	"strings"

	"github.com/goccy/go-json"
	"gopkg.in/yaml.v2"
)

func sanitize(filename string, content []byte) ([]byte, error) {
	if strings.HasSuffix(filename, ".yml") || strings.HasSuffix(filename, ".yml") {
		return sanitizeYaml(content)
	}
	if strings.HasSuffix(filename, ".json") {
		return sanitizeJson(content)
	}
	return content, fmt.Errorf("unsupported file type: %s", filename)
}

func sanitizeYaml(content []byte) ([]byte, error) {
	var data interface{}
	if err := yaml.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to decode as yaml: %w", err)
	}

	res, err := yaml.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode as yaml: %w", err)
	}

	return res, nil
}

func sanitizeJson(content []byte) ([]byte, error) {
	var data interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to decode as json: %w", err)
	}

	res, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to encode as json: %w", err)
	}

	return res, nil
}
