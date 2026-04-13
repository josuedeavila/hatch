package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// cmdInit scaffolds `.hatch/` with one example primitive of each kind so a
// new user can `hatch generate` immediately and see output.
func cmdInit(_ context.Context, args []string, stdout, stderr io.Writer) error {
	fs, root, _ := commonFlags("init", stderr)
	if err := fs.Parse(args); err != nil {
		return err
	}
	srcRoot := filepath.Join(*root, ".hatch")
	for _, sub := range []string{"rules", "skills", "commands", "agents"} {
		if err := os.MkdirAll(filepath.Join(srcRoot, sub), 0o755); err != nil {
			return err
		}
	}

	// Files in a slice so iteration order (and therefore stdout output) is
	// deterministic across runs. Map iteration in Go is randomized.
	type seed struct {
		path string
		body string
	}
	seeds := []seed{
		{
			path: filepath.Join(srcRoot, "rules", "coding-style.md"),
			body: "Write clean, well-tested Go. Prefer short functions. Use table-driven tests.\n",
		},
		{
			path: filepath.Join(srcRoot, "skills", "review-pr", "SKILL.md"),
			body: `---
description: Review an open pull request for correctness, style, and tests.
---

# review-pr

When asked to review a PR, check out the branch, read the diff end-to-end,
and report findings as: (1) bugs, (2) style nits, (3) missing tests.
`,
		},
		{
			path: filepath.Join(srcRoot, "commands", "commit.md"),
			body: `---
description: Stage and commit current changes with a generated message.
---

# commit

Summarise the staged diff in one sentence and create a commit with that
message.
`,
		},
		{
			path: filepath.Join(srcRoot, "agents", "security-auditor.md"),
			body: `---
description: Review code for common security pitfalls (injection, XSS, auth).
---

# security-auditor

Focus on OWASP Top 10 categories. Report findings as file:line references.
`,
		},
	}
	sort.Slice(seeds, func(i, j int) bool { return seeds[i].path < seeds[j].path })

	for _, s := range seeds {
		if _, err := os.Stat(s.path); err == nil {
			fmt.Fprintf(stdout, "exists  %s\n", s.path)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(s.path, []byte(s.body), 0o644); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "created %s\n", s.path)
	}
	return nil
}
