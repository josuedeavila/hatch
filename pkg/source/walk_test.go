package source_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/hatch/pkg/source"
	"github.com/matryer/is"
)

// writeTree creates a nested file structure under root from a map of
// relative paths to contents. Used to build test fixtures inline.
func writeTree(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for rel, body := range files {
		full := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestLoad_MissingSourceDirErrors(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	_, err := source.Load(dir)
	is.True(err != nil)
}

func TestLoad_Rules(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	writeTree(t, dir, map[string]string{
		".hatch/rules/alpha.md": "alpha body",
		".hatch/rules/beta.md":  "beta body",
	})
	src, err := source.Load(dir)
	is.NoErr(err)
	is.Equal(len(src.Rules), 2)
	// Sorted alphabetically by name.
	is.Equal(src.Rules[0].Name, "alpha")
	is.Equal(src.Rules[1].Name, "beta")
}

func TestLoad_SkillsAsDirectories(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	writeTree(t, dir, map[string]string{
		".hatch/skills/review-pr/SKILL.md": "---\ndescription: review\n---\nbody\n",
	})
	src, err := source.Load(dir)
	is.NoErr(err)
	is.Equal(len(src.Skills), 1)
	is.Equal(src.Skills[0].Name, "review-pr")
	is.Equal(src.Skills[0].Description, "review")
	// SourcePath points to the skill directory, not the SKILL.md file.
	info, err := os.Stat(src.Skills[0].SourcePath)
	is.NoErr(err)
	is.True(info.IsDir())
}

func TestLoad_CommandsAndAgents(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	writeTree(t, dir, map[string]string{
		".hatch/commands/commit.md": "---\ndescription: c\n---\nbody\n",
		".hatch/agents/security.md": "---\ndescription: a\n---\nbody\n",
	})
	src, err := source.Load(dir)
	is.NoErr(err)
	is.Equal(len(src.Commands), 1)
	is.Equal(len(src.Agents), 1)
	is.Equal(src.Commands[0].Name, "commit")
	is.Equal(src.Agents[0].Name, "security")
}

func TestLoad_IgnoresNonMarkdownFiles(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	writeTree(t, dir, map[string]string{
		".hatch/rules/real.md":   "content",
		".hatch/rules/notes.txt": "ignore me",
	})
	src, err := source.Load(dir)
	is.NoErr(err)
	is.Equal(len(src.Rules), 1)
}
