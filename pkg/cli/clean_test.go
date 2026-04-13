package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestClean_RemovesFilesAndStripsBlocks(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)

	_, _, err := invoke(t, "generate", "-C", dir)
	is.NoErr(err)

	_, _, err = invoke(t, "clean", "-C", dir)
	is.NoErr(err)

	_, err = os.Stat(filepath.Join(dir, ".claude/skills/review-pr/SKILL.md"))
	is.True(os.IsNotExist(err))

	// With only hatch content, CLAUDE.md is removed entirely.
	_, err = os.Stat(filepath.Join(dir, "CLAUDE.md"))
	is.True(os.IsNotExist(err))
}

func TestClean_PreservesSurroundingUserContent(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)

	userFile := filepath.Join(dir, "CLAUDE.md")
	is.NoErr(os.WriteFile(userFile, []byte("# User\n\nKeep me.\n"), 0o644))

	_, _, err := invoke(t, "generate", "-C", dir)
	is.NoErr(err)
	_, _, err = invoke(t, "clean", "-C", dir)
	is.NoErr(err)

	got, err := os.ReadFile(userFile)
	is.NoErr(err)
	is.True(strings.Contains(string(got), "# User"))
	is.True(strings.Contains(string(got), "Keep me."))
	is.True(!strings.Contains(string(got), "hatch:begin"))
}
