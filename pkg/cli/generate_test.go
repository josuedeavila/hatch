package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestGenerate_WritesExpectedFilesAcrossTargets(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)

	stdout, _, err := invoke(t, "generate", "-C", dir)
	is.NoErr(err)
	is.True(strings.Contains(stdout, "CLAUDE.md"))
	is.True(strings.Contains(stdout, "AGENTS.md"))
	is.True(strings.Contains(stdout, ".github/copilot-instructions.md"))

	claudeMD, err := os.ReadFile(filepath.Join(dir, "CLAUDE.md"))
	is.NoErr(err)
	is.True(strings.Contains(string(claudeMD), "<!-- hatch:begin v1 -->"))
	is.True(strings.Contains(string(claudeMD), "Be concise."))

	skill, err := os.ReadFile(filepath.Join(dir, ".claude/skills/review-pr/SKILL.md"))
	is.NoErr(err)
	is.True(strings.Contains(string(skill), "description: Review a PR."))
}

func TestGenerate_Idempotent(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)

	_, _, err := invoke(t, "generate", "-C", dir)
	is.NoErr(err)
	first, err := os.ReadFile(filepath.Join(dir, "CLAUDE.md"))
	is.NoErr(err)

	_, _, err = invoke(t, "generate", "-C", dir)
	is.NoErr(err)
	second, err := os.ReadFile(filepath.Join(dir, "CLAUDE.md"))
	is.NoErr(err)

	is.Equal(string(first), string(second))
}

func TestGenerate_PreservesUserContentAroundBlock(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)

	userFile := filepath.Join(dir, "CLAUDE.md")
	userContent := "# My notes\n\nHand-written content here.\n"
	is.NoErr(os.WriteFile(userFile, []byte(userContent), 0o644))

	_, _, err := invoke(t, "generate", "-C", dir)
	is.NoErr(err)

	got, err := os.ReadFile(userFile)
	is.NoErr(err)
	s := string(got)
	is.True(strings.Contains(s, "# My notes"))
	is.True(strings.Contains(s, "Hand-written content here."))
	is.True(strings.Contains(s, "<!-- hatch:begin v1 -->"))
	is.True(strings.Contains(s, "Be concise."))
}

func TestGenerate_TargetsFlag(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)

	stdout, _, err := invoke(t, "generate", "-C", dir, "-targets", "claude")
	is.NoErr(err)
	is.True(strings.Contains(stdout, "CLAUDE.md"))
	is.True(!strings.Contains(stdout, "AGENTS.md"))
}

func TestGenerate_UnknownTargetErrors(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)
	_, _, err := invoke(t, "generate", "-C", dir, "-targets", "nosuch")
	is.True(err != nil)
}

func TestGenerate_AliasGen(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)
	_, _, err := invoke(t, "gen", "-C", dir)
	is.NoErr(err)
}
