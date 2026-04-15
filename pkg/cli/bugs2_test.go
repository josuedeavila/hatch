package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
)

// Bug 6: Codex and OpenCode both write to AGENTS.md. If a rule filters to
// one of them (e.g., `targets: [codex]`), the target that writes LAST
// silently overwrites the earlier write with its own (filtered) view of
// the rules, losing the filtered-in rule. High severity (silent data loss:
// content that the user asked to appear disappears).
func TestBug_AGENTSMdSharedByCodexAndOpenCode(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	t.Chdir(dir)

	// Two rules: one shared by all, one targeting codex only. The
	// codex-only rule must still appear in AGENTS.md — OpenCode's write
	// of the same file must not erase it.
	writeFile(t, ".hatch/_rules/shared.md", "SHARED RULE CONTENT\n")
	writeFile(t, ".hatch/_rules/codex-only.md", `---
targets: [codex]
---
CODEX ONLY RULE CONTENT
`)

	_, _, err := invoke(t, "gen")
	is.NoErr(err)

	got, err := os.ReadFile("AGENTS.md")
	is.NoErr(err)
	s := string(got)
	is.True(strings.Contains(s, "SHARED RULE CONTENT"))
	is.True(strings.Contains(s, "CODEX ONLY RULE CONTENT"))
}

func writeFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}
