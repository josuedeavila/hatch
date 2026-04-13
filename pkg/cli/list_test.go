package cli_test

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestList_PrintsTargetsAndArtifacts(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	scaffoldSource(t, dir)
	t.Chdir(dir)

	stdout, _, err := invoke(t, "list")
	is.NoErr(err)
	is.True(strings.Contains(stdout, "Claude Code (claude)"))
	is.True(strings.Contains(stdout, "CLAUDE.md"))
	is.True(strings.Contains(stdout, "[block]"))
	is.True(strings.Contains(stdout, "[file]"))
}

func TestList_NoSourceErrors(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)
	_, _, err := invoke(t, "list")
	is.True(err != nil)
}
