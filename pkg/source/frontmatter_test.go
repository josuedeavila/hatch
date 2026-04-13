package source

import (
	"testing"

	"github.com/matryer/is"
)

func TestSplitFrontmatter(t *testing.T) {
	cases := []struct {
		name     string
		in       string
		wantFM   string
		wantBody string
		wantErr  bool
	}{
		{
			name:     "no frontmatter",
			in:       "# hello\nworld\n",
			wantFM:   "",
			wantBody: "# hello\nworld\n",
		},
		{
			name:     "simple frontmatter",
			in:       "---\nname: x\n---\nbody\n",
			wantFM:   "name: x",
			wantBody: "body\n",
		},
		{
			name:     "frontmatter with leading whitespace",
			in:       "\n\n---\nname: x\n---\nbody\n",
			wantFM:   "name: x",
			wantBody: "body\n",
		},
		{
			name:    "unterminated frontmatter",
			in:      "---\nname: x\nno end here",
			wantErr: true,
		},
		{
			name:     "empty body",
			in:       "---\nname: x\n---\n",
			wantFM:   "name: x",
			wantBody: "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			fm, body, err := splitFrontmatter([]byte(tc.in))
			if tc.wantErr {
				is.True(err != nil) // expected error
				return
			}
			is.NoErr(err)
			is.Equal(string(fm), tc.wantFM)
			is.Equal(string(body), tc.wantBody)
		})
	}
}

func TestParsePrimitive(t *testing.T) {
	t.Run("skill with required fields", func(t *testing.T) {
		is := is.New(t)
		data := []byte(`---
name: review-pr
description: Review a PR end-to-end.
---
body text
`)
		p, err := parsePrimitive(KindSkill, "review-pr", "src.md", data)
		is.NoErr(err)
		is.Equal(p.Name, "review-pr")
		is.Equal(p.Description, "Review a PR end-to-end.")
		is.Equal(p.Body, "body text\n")
	})

	t.Run("skill missing description errors", func(t *testing.T) {
		is := is.New(t)
		data := []byte(`---
name: review-pr
---
body
`)
		_, err := parsePrimitive(KindSkill, "review-pr", "src.md", data)
		is.True(err != nil) // missing description should error
	})

	t.Run("rule without frontmatter is valid", func(t *testing.T) {
		is := is.New(t)
		data := []byte("# coding style\nbe concise\n")
		p, err := parsePrimitive(KindRule, "coding-style", "rule.md", data)
		is.NoErr(err)
		is.Equal(p.Name, "coding-style")
		is.Equal(p.Body, "# coding style\nbe concise\n")
	})

	t.Run("rule with applyTo", func(t *testing.T) {
		is := is.New(t)
		data := []byte(`---
applyTo: "**/*.go"
---
Go rules here
`)
		p, err := parsePrimitive(KindRule, "go-rules", "r.md", data)
		is.NoErr(err)
		is.Equal(p.ApplyTo, "**/*.go")
	})

	t.Run("targets list", func(t *testing.T) {
		is := is.New(t)
		data := []byte(`---
description: something
targets:
  - claude
  - opencode
---
body
`)
		p, err := parsePrimitive(KindSkill, "thing", "s.md", data)
		is.NoErr(err)
		is.Equal(len(p.Targets), 2)
		is.Equal(p.Targets[0], "claude")
		is.Equal(p.Targets[1], "opencode")
	})

	t.Run("per-target passthrough captured in Overrides", func(t *testing.T) {
		is := is.New(t)
		data := []byte(`---
description: something
claude:
  allowed-tools:
    - Read
    - Grep
---
body
`)
		p, err := parsePrimitive(KindSkill, "thing", "s.md", data)
		is.NoErr(err)
		is.True(p.Overrides != nil)
		_, ok := p.Overrides["claude"]
		is.True(ok) // claude override should be present
	})
}
