package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/matryer/hatch/pkg/block"
	"github.com/matryer/hatch/pkg/source"
	"github.com/matryer/hatch/pkg/target"
)

// cmdClean computes what a fresh build WOULD write, then removes those files
// (for ModeFile) or strips just the hatch block (for ModeBlock). No manifest
// is kept: the source tree is the single source of truth for what exists.
func cmdClean(_ context.Context, available *target.Set, args []string, stdout, stderr io.Writer) error {
	fs, root, targetsList := commonFlags("clean", stderr)
	if err := fs.Parse(args); err != nil {
		return err
	}
	targets, err := selectTargets(available, *targetsList)
	if err != nil {
		return err
	}
	src, err := source.Load(*root)
	if err != nil {
		return err
	}

	for _, t := range targets.All() {
		arts, err := t.Emit(src)
		if err != nil {
			return fmt.Errorf("%s: %w", t.Name(), err)
		}
		for _, a := range arts {
			if err := cleanArtifact(*root, a); err != nil {
				fmt.Fprintf(stderr, "hatch: %s: %s\n", a.Path, err)
				continue
			}
			fmt.Fprintf(stdout, "removed %s (%s)\n", a.Path, a.Mode)
		}
	}
	return nil
}

func cleanArtifact(root string, a target.Artifact) error {
	full := filepath.Join(root, a.Path)
	switch a.Mode {
	case target.ModeFile:
		if err := os.Remove(full); err != nil && !os.IsNotExist(err) {
			return err
		}
		pruneEmptyDirs(root, filepath.Dir(full))
		return nil
	case target.ModeBlock:
		return block.Strip(full, block.CurrentMarker)
	default:
		return fmt.Errorf("unknown write mode %q", a.Mode)
	}
}

// pruneEmptyDirs walks upward from dir, removing directories that become
// empty, stopping at root.
func pruneEmptyDirs(root, dir string) {
	absRoot, _ := filepath.Abs(root)
	for {
		absDir, _ := filepath.Abs(dir)
		if absDir == absRoot || len(absDir) <= len(absRoot) {
			return
		}
		entries, err := os.ReadDir(dir)
		if err != nil || len(entries) > 0 {
			return
		}
		if err := os.Remove(dir); err != nil {
			return
		}
		dir = filepath.Dir(dir)
	}
}
