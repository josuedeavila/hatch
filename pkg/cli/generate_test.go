package cli_test

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestGen_WritesExpectedFilesAcrossTargets(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)
	t.Chdir(dir)

	stdout, _, err := invoke(t, "gen")
	is.NoErr(err)
	is.True(strings.Contains(stdout, "CLAUDE.md"))
	is.True(strings.Contains(stdout, "AGENTS.md"))
	is.True(strings.Contains(stdout, ".github/copilot-instructions.md"))

	claudeMD, err := os.ReadFile("CLAUDE.md")
	is.NoErr(err)
	is.True(strings.Contains(string(claudeMD), "<!-- hatch:begin v1 -->"))
	is.True(strings.Contains(string(claudeMD), "Be concise."))

	skill, err := os.ReadFile(".claude/skills/review-pr/SKILL.md")
	is.NoErr(err)
	is.True(strings.Contains(string(skill), "description: Review a PR."))
}

func TestGen_Idempotent(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)
	t.Chdir(dir)

	_, _, err := invoke(t, "gen")
	is.NoErr(err)
	first, err := os.ReadFile("CLAUDE.md")
	is.NoErr(err)

	_, _, err = invoke(t, "gen")
	is.NoErr(err)
	second, err := os.ReadFile("CLAUDE.md")
	is.NoErr(err)

	is.Equal(string(first), string(second))
}

func TestGen_PreservesUserContentAroundBlock(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)
	t.Chdir(dir)

	userContent := "# My notes\n\nHand-written content here.\n"
	is.NoErr(os.WriteFile("CLAUDE.md", []byte(userContent), 0o644))

	_, _, err := invoke(t, "gen")
	is.NoErr(err)

	got, err := os.ReadFile("CLAUDE.md")
	is.NoErr(err)
	s := string(got)
	is.True(strings.Contains(s, "# My notes"))
	is.True(strings.Contains(s, "Hand-written content here."))
	is.True(strings.Contains(s, "<!-- hatch:begin v1 -->"))
	is.True(strings.Contains(s, "Be concise."))
}

func TestGen_TargetsFlag(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)
	t.Chdir(dir)

	stdout, _, err := invoke(t, "gen", "-targets", "claude")
	is.NoErr(err)
	is.True(strings.Contains(stdout, "CLAUDE.md"))
	is.True(!strings.Contains(stdout, "AGENTS.md"))
}

func TestGen_UnknownTargetErrors(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)
	t.Chdir(dir)
	_, _, err := invoke(t, "gen", "-targets", "nosuch")
	is.True(err != nil)
}

func TestGen_UnknownTargetInListWritesNothing(t *testing.T) {
	// When the -targets list mixes a valid target with an unknown one,
	// the command must fail before touching the filesystem. Nothing is
	// written — not even the valid target's output.
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)
	t.Chdir(dir)

	_, _, err := invoke(t, "gen", "-targets", "claude,nosuch")
	is.True(err != nil)
	is.True(strings.Contains(err.Error(), "nosuch"))

	// Claude's file must NOT exist — the error should fire before any
	// target-level writes.
	for _, path := range []string{
		"CLAUDE.md",
		".claude/skills/review-pr/SKILL.md",
		"AGENTS.md",
	} {
		_, statErr := os.Stat(path)
		is.True(os.IsNotExist(statErr))
	}
}

func TestGen_PositionalArgErrors(t *testing.T) {
	// Target selection is -targets only. A stray positional word must
	// fail loudly rather than be silently ignored and run against every
	// target.
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)
	t.Chdir(dir)

	_, _, err := invoke(t, "gen", "claude")
	is.True(err != nil)
	is.True(strings.Contains(err.Error(), "unexpected argument"))

	// No files should have been written.
	for _, path := range []string{"CLAUDE.md", "AGENTS.md", ".claude/skills/review-pr/SKILL.md"} {
		_, statErr := os.Stat(path)
		is.True(os.IsNotExist(statErr))
	}
}

func TestGen_RespectsCanceledContext(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)
	t.Chdir(dir)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, _, err := invokeWithCtx(t, ctx, "", "gen")
	is.True(err != nil)
	is.True(errors.Is(err, context.Canceled))

	// Nothing should have been written.
	_, err = os.Stat("CLAUDE.md")
	is.True(os.IsNotExist(err))
}

func TestGen_LegacyGenerateWordRemoved(t *testing.T) {
	// `hatch generate` (the old long form) should now be an unknown command.
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)
	t.Chdir(dir)
	_, _, err := invoke(t, "generate")
	is.True(err != nil)
}
