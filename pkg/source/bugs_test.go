package source

import (
	"testing"

	"github.com/matryer/is"
)

// Bug 1: A scalar `targets:` value (single string, not a YAML list) is
// silently dropped — HasTarget would then return "all targets" instead of
// the single target the user intended. High severity (silent data loss).
func TestBug_ScalarTargetsShouldBeSingleton(t *testing.T) {
	is := is.New(t)
	data := []byte(`---
description: something
targets: claude
---
body
`)
	p, err := parsePrimitive(KindSkill, "thing", "s.md", data)
	is.NoErr(err)
	is.Equal(len(p.Targets), 1)
	is.Equal(p.Targets[0], "claude")
}

// Bug 2: Frontmatter with CRLF line endings parses but pollutes the body
// with \r characters, producing mixed-line-ending generated files. Medium
// severity: cosmetic and cross-platform inconsistency.
func TestBug_CRLFFrontmatterNormalisesBody(t *testing.T) {
	is := is.New(t)
	data := []byte("---\r\nname: review-pr\r\ndescription: r\r\n---\r\nline one\r\nline two\r\n")
	p, err := parsePrimitive(KindSkill, "review-pr", "s.md", data)
	is.NoErr(err)
	is.Equal(p.Name, "review-pr")
	is.Equal(p.Description, "r")
	// Body should be LF-normalised — no stray CR characters leaking through.
	for i := 0; i < len(p.Body); i++ {
		if p.Body[i] == '\r' {
			t.Fatalf("body contains CR byte at index %d: %q", i, p.Body)
		}
	}
}

// Bug 3: An empty `name:` value in frontmatter wipes the name derived from
// the filename, producing a confusing "missing name" validation error.
// Medium severity (bad UX).
func TestBug_EmptyNameInFrontmatterKeepsDerived(t *testing.T) {
	is := is.New(t)
	data := []byte(`---
name:
description: thing
---
body
`)
	p, err := parsePrimitive(KindSkill, "derived-name", "s.md", data)
	is.NoErr(err)
	is.Equal(p.Name, "derived-name")
}
