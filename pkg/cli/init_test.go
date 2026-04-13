package cli_test

import (
	"os"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestInit_ScaffoldsAllPrimitives(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)

	stdout, _, err := invoke(t, "init")
	is.NoErr(err)
	is.True(strings.Contains(stdout, "created"))

	for _, rel := range []string{
		".hatch/rules/coding-style.md",
		".hatch/skills/review-pr/SKILL.md",
		".hatch/commands/commit.md",
		".hatch/agents/security-auditor.md",
	} {
		_, err := os.Stat(rel)
		is.NoErr(err) // scaffold file should exist
	}
}

func TestInit_SkipsExistingFiles(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)

	// First init creates files.
	_, _, err := invoke(t, "init")
	is.NoErr(err)

	// Overwrite a file with custom content.
	rulePath := ".hatch/rules/coding-style.md"
	is.NoErr(os.WriteFile(rulePath, []byte("custom content"), 0o644))

	// Second init should NOT overwrite the user's custom content.
	_, _, err = invoke(t, "init")
	is.NoErr(err)
	got, err := os.ReadFile(rulePath)
	is.NoErr(err)
	is.Equal(string(got), "custom content")
}
