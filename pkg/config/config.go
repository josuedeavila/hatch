// Package config loads `.hatch/config.yaml`. If the file is absent, all known
// targets are enabled and output is written to the project root.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config is the parsed `.hatch/config.yaml`.
type Config struct {
	Targets []string `yaml:"targets"`
	Output  string   `yaml:"output"`
}

// DefaultTargets is the full set of targets hatch knows how to generate.
var DefaultTargets = []string{"claude", "codex", "copilot", "opencode"}

// Load reads `<root>/.hatch/config.yaml`. A missing file is not an error — the
// defaults are returned.
func Load(root string) (*Config, error) {
	cfg := &Config{
		Targets: append([]string(nil), DefaultTargets...),
		Output:  ".",
	}
	path := filepath.Join(root, ".hatch", "config.yaml")
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if cfg.Output == "" {
		cfg.Output = "."
	}
	if len(cfg.Targets) == 0 {
		cfg.Targets = append([]string(nil), DefaultTargets...)
	}
	return cfg, nil
}
