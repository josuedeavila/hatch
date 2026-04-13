package render

import (
	"testing"

	"github.com/matryer/is"
	"gopkg.in/yaml.v3"
)

func TestFrontmatter_Empty(t *testing.T) {
	is := is.New(t)
	out, err := Frontmatter(nil)
	is.NoErr(err)
	is.Equal(out, "")
}

func TestFrontmatter_SimpleFields(t *testing.T) {
	is := is.New(t)
	out, err := Frontmatter([]Field{
		{Key: "name", Value: "review-pr"},
		{Key: "description", Value: "Review a PR."},
	})
	is.NoErr(err)
	want := "---\nname: review-pr\ndescription: Review a PR.\n---\n"
	is.Equal(out, want)
}

func TestFrontmatter_SkipsEmptyValues(t *testing.T) {
	is := is.New(t)
	out, err := Frontmatter([]Field{
		{Key: "name", Value: "x"},
		{Key: "description", Value: ""},
		{Key: "applyTo", Value: ""},
	})
	is.NoErr(err)
	is.Equal(out, "---\nname: x\n---\n")
}

func TestFrontmatter_StringSlice(t *testing.T) {
	is := is.New(t)
	out, err := Frontmatter([]Field{
		{Key: "targets", Value: []string{"claude", "opencode"}},
	})
	is.NoErr(err)
	// Flow style is used for short string slices.
	is.Equal(out, "---\ntargets: [claude, opencode]\n---\n")
}

func TestFrontmatter_Deterministic(t *testing.T) {
	// Running Frontmatter twice with the same input must produce byte-identical
	// output (no map iteration order leakage).
	is := is.New(t)
	fields := []Field{
		{Key: "name", Value: "x"},
		{Key: "description", Value: "y"},
		{Key: "targets", Value: []string{"a", "b", "c"}},
	}
	a, err := Frontmatter(fields)
	is.NoErr(err)
	b, err := Frontmatter(fields)
	is.NoErr(err)
	is.Equal(a, b)
}

func TestMergeOverride_ReplacesKnown(t *testing.T) {
	is := is.New(t)
	fields := []Field{
		{Key: "name", Value: "x"},
		{Key: "description", Value: "old"},
	}
	var n yaml.Node
	err := yaml.Unmarshal([]byte("description: new\nmodel: claude-opus-4\n"), &n)
	is.NoErr(err)
	override := n.Content[0]
	merged := MergeOverride(fields, override)
	is.Equal(len(merged), 3)
	// "description" replaced in place; "model" appended.
	is.Equal(merged[0].Key, "name")
	is.Equal(merged[1].Key, "description")
	is.Equal(merged[2].Key, "model")
}
