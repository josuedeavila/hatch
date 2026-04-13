package cli_test

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestClean_RemovesFilesAndStripsBlocks(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)
	t.Chdir(dir)

	_, _, err := invoke(t, "gen")
	is.NoErr(err)

	_, _, err = invoke(t, "clean")
	is.NoErr(err)

	_, err = os.Stat(".claude/skills/review-pr/SKILL.md")
	is.True(os.IsNotExist(err))

	// With only hatch content, CLAUDE.md is removed entirely.
	_, err = os.Stat("CLAUDE.md")
	is.True(os.IsNotExist(err))
}

func TestClean_RespectsCanceledContext(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)
	t.Chdir(dir)

	// First, generate so there's something to clean.
	_, _, err := invoke(t, "gen")
	is.NoErr(err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, _, err = invokeWithCtx(t, ctx, "", "clean")
	is.True(err != nil)
	is.True(errors.Is(err, context.Canceled))

	// Generated files should still exist — clean bailed before touching them.
	_, err = os.Stat("CLAUDE.md")
	is.NoErr(err)
}

func TestClean_PreservesSurroundingUserContent(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)
	t.Chdir(dir)

	is.NoErr(os.WriteFile("CLAUDE.md", []byte("# User\n\nKeep me.\n"), 0o644))

	_, _, err := invoke(t, "gen")
	is.NoErr(err)
	_, _, err = invoke(t, "clean")
	is.NoErr(err)

	got, err := os.ReadFile("CLAUDE.md")
	is.NoErr(err)
	is.True(strings.Contains(string(got), "# User"))
	is.True(strings.Contains(string(got), "Keep me."))
	is.True(!strings.Contains(string(got), "hatch:begin"))
}
