package cli_test

import (
	"context"
	"errors"
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

func TestClean_UnknownTargetInListDoesNothing(t *testing.T) {
	// Mixed valid+invalid targets → whole command fails before touching
	// anything. Generated files from a prior `hatch gen` must survive.
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)
	t.Chdir(dir)

	_, _, err := invoke(t, "gen")
	is.NoErr(err)

	_, _, err = invoke(t, "clean", "-targets", "claude,nosuch")
	is.True(err != nil)
	is.True(strings.Contains(err.Error(), "nosuch"))

	// Claude's skill file must still be there — clean bailed before
	// removing anything.
	_, statErr := os.Stat(".claude/skills/review-pr/SKILL.md")
	is.NoErr(statErr)
}

func TestClean_PositionalArgErrors(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)
	t.Chdir(dir)
	_, _, err := invoke(t, "clean", "claude")
	is.True(err != nil)
	is.True(strings.Contains(err.Error(), "unexpected argument"))
}

func TestClean_NestedPath_RemovesAndPrunes(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	files := map[string]string{
		".hatch/_rules/global.md":               "GLOBAL\n",
		".hatch/backend/_rules/api.md":          "BACKEND\n",
		".hatch/backend/_skills/check/SKILL.md": "---\ndescription: c\n---\nbody\n",
	}
	for path, body := range files {
		full := dir + "/" + path
		is.NoErr(os.MkdirAll(filepath.Dir(full), 0o755))
		is.NoErr(os.WriteFile(full, []byte(body), 0o644))
	}
	t.Chdir(dir)

	_, _, err := invoke(t, "gen")
	is.NoErr(err)
	// Sanity: the nested files exist before clean.
	_, err = os.Stat("backend/CLAUDE.md")
	is.NoErr(err)
	_, err = os.Stat("backend/.claude/skills/check/SKILL.md")
	is.NoErr(err)

	_, _, err = invoke(t, "clean")
	is.NoErr(err)

	// All hatch-owned nested files gone.
	for _, p := range []string{
		"backend/CLAUDE.md",
		"backend/AGENTS.md",
		"backend/.claude/skills/check/SKILL.md",
		"backend/.claude",
		"backend/.opencode",
		"backend/.agents",
	} {
		_, err := os.Stat(p)
		is.True(os.IsNotExist(err)) // hatch-owned nested file/dir must be gone
	}
	// And the .hatch/ source tree is untouched.
	_, err = os.Stat(".hatch/backend/_rules/api.md")
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
