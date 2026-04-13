package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

func TestLoad_Missing_ReturnsDefaults(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	cfg, err := Load(dir)
	is.NoErr(err)
	is.Equal(cfg.Output, ".")
	is.Equal(len(cfg.Targets), len(DefaultTargets))
}

func TestLoad_NarrowsTargets(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".hatch")
	is.NoErr(os.MkdirAll(cfgDir, 0o755))
	is.NoErr(os.WriteFile(filepath.Join(cfgDir, "config.yaml"), []byte("targets: [claude, codex]\n"), 0o644))

	cfg, err := Load(dir)
	is.NoErr(err)
	is.Equal(len(cfg.Targets), 2)
	is.Equal(cfg.Targets[0], "claude")
	is.Equal(cfg.Targets[1], "codex")
}

func TestLoad_EmptyTargetsFallsBackToDefault(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".hatch")
	is.NoErr(os.MkdirAll(cfgDir, 0o755))
	is.NoErr(os.WriteFile(filepath.Join(cfgDir, "config.yaml"), []byte("output: out\n"), 0o644))

	cfg, err := Load(dir)
	is.NoErr(err)
	is.Equal(cfg.Output, "out")
	is.Equal(len(cfg.Targets), len(DefaultTargets))
}
