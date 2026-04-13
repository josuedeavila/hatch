package block

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestRender(t *testing.T) {
	is := is.New(t)
	out := Render(CurrentMarker, "hello world")
	want := "<!-- hatch:begin v1 -->\nhello world\n<!-- hatch:end v1 -->\n"
	is.Equal(out, want)
}

func TestInject_NewFile(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "CLAUDE.md")
	err := Inject(path, CurrentMarker, "rules body")
	is.NoErr(err)
	got, err := os.ReadFile(path)
	is.NoErr(err)
	is.True(strings.Contains(string(got), "<!-- hatch:begin v1 -->"))
	is.True(strings.Contains(string(got), "rules body"))
	is.True(strings.Contains(string(got), "<!-- hatch:end v1 -->"))
}

func TestInject_ReplacesExistingBlock(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "CLAUDE.md")
	initial := "# User content\n\n<!-- hatch:begin v1 -->\nold body\n<!-- hatch:end v1 -->\n\n# More user content\n"
	is.NoErr(os.WriteFile(path, []byte(initial), 0o644))

	err := Inject(path, CurrentMarker, "new body")
	is.NoErr(err)
	got, err := os.ReadFile(path)
	is.NoErr(err)
	s := string(got)
	is.True(strings.Contains(s, "# User content"))          // preserved before
	is.True(strings.Contains(s, "# More user content"))     // preserved after
	is.True(strings.Contains(s, "new body"))                // replaced
	is.True(!strings.Contains(s, "old body"))               // old gone
	is.True(strings.Contains(s, "<!-- hatch:begin v1 -->")) // markers intact
}

func TestInject_AppendsToFileWithoutMarkers(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "CLAUDE.md")
	is.NoErr(os.WriteFile(path, []byte("# User content\n"), 0o644))

	err := Inject(path, CurrentMarker, "hatch body")
	is.NoErr(err)
	got, _ := os.ReadFile(path)
	s := string(got)
	is.True(strings.Contains(s, "# User content")) // preserved
	is.True(strings.Contains(s, "hatch body"))     // appended
}

func TestInject_Idempotent(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "CLAUDE.md")

	// First injection into empty space.
	is.NoErr(Inject(path, CurrentMarker, "same body"))
	first, _ := os.ReadFile(path)

	// Second injection with identical content must produce identical output.
	is.NoErr(Inject(path, CurrentMarker, "same body"))
	second, _ := os.ReadFile(path)

	is.Equal(string(first), string(second))
}

func TestStrip_RemovesBlockKeepsNeighbors(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "CLAUDE.md")
	content := "# header\n\n<!-- hatch:begin v1 -->\nhatch content\n<!-- hatch:end v1 -->\n\n# footer\n"
	is.NoErr(os.WriteFile(path, []byte(content), 0o644))

	is.NoErr(Strip(path, CurrentMarker))
	got, _ := os.ReadFile(path)
	s := string(got)
	is.True(strings.Contains(s, "# header"))
	is.True(strings.Contains(s, "# footer"))
	is.True(!strings.Contains(s, "hatch content"))
	is.True(!strings.Contains(s, "hatch:begin"))
}

func TestStrip_RemovesFileWhenOnlyBlock(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "CLAUDE.md")
	is.NoErr(Inject(path, CurrentMarker, "just the block"))

	is.NoErr(Strip(path, CurrentMarker))
	_, err := os.Stat(path)
	is.True(os.IsNotExist(err)) // file should be removed
}
